package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/digitalrebar/provision/client/events"
	"github.com/digitalrebar/provision/client/jobs"
	models "github.com/digitalrebar/provision/genmodels"
	"github.com/go-openapi/strfmt"
	"github.com/spf13/cobra"
)

var exitOnFailure = false
var runnerDir string
var actuallyPowerThings = true

var cmdHelper = []byte(`
#!/bin/bash

# To force dpkg on Debian-based distros to play nice.
export DEBIAN_FRONTEND=noninteractive DEBCONF_NONINTERACTIVE_SEEN=true

# Force everything to use the C locale to keep things sane
export LC_ALL=C LANGUAGE=C LANG=C

# Make sure we play nice with debugging
export PS4='${BASH_SOURCE}@${LINENO}(${FUNCNAME[0]}): '

# Make sure the scripts are somewhat typo-resistant
set -o pipefail -o errexit
shopt -s nullglob extglob globstar

# Make sure that $PATH is somewhat sane.
fix_path() {
    local -A pathparts
    local part
    local IFS=':'
    for part in $PATH; do
        pathparts["$part"]="true"
    done
    local wanted_pathparts=("/usr/local/bin" "/usr/local/sbin" "/bin" "/sbin" "/usr/bin" "/usr/sbin")
    for part in "${wanted_pathparts[@]}"; do
        [[ ${pathparts[$part]} ]] && continue
        PATH="$part:$PATH"
    done
}
fix_path
unset fix_path

# Figure out what Linux distro we are running on.
export OS_TYPE= OS_VER= OS_NAME=
if [[ -f /etc/os-release ]]; then
    . /etc/os-release
    OS_TYPE=${ID,,}
    OS_VER=${VERSION_ID,,}
elif [[ -f /etc/lsb-release ]]; then
    . /etc/lsb-release
    OS_VER=${DISTRIB_RELEASE,,}
    OS_TYPE=${DISTRIB_ID,,}
elif [[ -f /etc/centos-release || -f /etc/fedora-release || -f /etc/redhat-release ]]; then
    for rel in centos-release fedora-release redhat-release; do
        [[ -f /etc/$rel ]] || continue
        OS_TYPE=${rel%%-*}
        OS_VER="$(egrep -o '[0-9.]+' "/etc/$rel")"
        break
    done
    if [[ ! $OS_TYPE ]]; then
        echo "Cannot determine Linux version we are running on!"
        exit 1
    fi
elif [[ -f /etc/debian_version ]]; then
    OS_TYPE=debian
    OS_VER=$(cat /etc/debian_version)
fi
OS_NAME="$OS_TYPE-$OS_VER"

case $OS_TYPE in
    centos|redhat|fedora|rhel|scientificlinux) OS_FAMILY="rhel";;
    debian|ubuntu) OS_FAMILY="debian";;
    *) OS_FAMILY=$OS_TYPE;;
esac

if_update_needed() {
    local timestampref=/tmp/pkg_cache_update
    if [[ ! -f $timestampref ]] || \
           (( ($(stat -c '%Y' "$timestampref") - $(date '+%s')) > 86400 )); then
        touch "$timestampref"
        "$@"
    fi
}

# Install a package
install() {
    local to_install=()
    local pkg
    for pkg in "$@"; do
        to_install+=("$pkg")
    done
    case $OS_FAMILY in
        rhel)
            if_update_needed yum -y makecache
            yum -y install "${to_install[@]}";;
        debian)
            if_update_needed apt-get -y update
            apt-get -y install "${to_install[@]}";;
        alpine)
            if_update_needed apk update
            apk add "${to_install[@]}";;
        *) echo "No idea how to install packages on $OS_NAME"
           exit 1;;
    esac
}

INITSTYLE="sysv"
if which systemctl &>/dev/null; then
    INITSTYLE="systemd"
elif which initctl &>/dev/null; then
    INITSTYLE="upstart"
fi

# Perform service actions.
service() {
    # $1 = service name
    # $2 = action to perform
    local svc="$1"
    shift
    if which systemctl &>/dev/null; then
        systemctl "$1" "$svc.service"
    elif which chkconfig &>/dev/null; then
        case $1 in
            enable) chkconfig "$svc" on;;
            disable) chkconfig "$svc" off;;
            *)  command service "$svc" "$@";;
        esac
    elif which initctl &>/dev/null && initctl version 2>/dev/null | grep -q upstart ; then
        /usr/sbin/service "$svc" "$1"
    elif [[ -f /etc/init/$svc.unit ]]; then
        initctl "$1" "$svc"
    elif which update-rc.d &>/dev/null; then
        case $1 in
            enable|disable) update-rc.d "$svc" "$1";;
            *) "/etc/init.d/$svc" "$1";;
        esac
    elif [[ -x /etc/init.d/$svc ]]; then
        "/etc/init.d/$svc" "$1"
    else
        echo "No idea how to manage services on $OS_NAME"
        exit 1
    fi
}

get_param() {
    # $1 attrib to get.  Attrib will be fetched in the context of the current machine
    local attr
    drpcli machines get "$RS_UUID" param "$1"
}

set_param() {
    # $1 = name of the parameter to set
    # $2 = parameter to set.
    #      if $2 == "", then we will read from stdin
    local src="$2"
    if [[ ! $src ]]; then src="-"; fi
    drpcli machines set "$RS_UUID" param "$1" to "$src"
}

__sane_exit() {
    touch "$RS_TASK_DIR/.sane-exit-codes"
}

__exit() {
    __sane_exit
    exit $1
}

exit_incomplete() {
    __exit 128
}

exit_reboot() {
    __exit 64
}

exit_shutdown() {
    __exit 32
}

exit_incomplete_reboot() {
    __exit 192
}

exit_incomplete_shutdown() {
    __exit 160
}

addr_port() {
    if [[ $1 =~ ':' ]]; then
        printf '[%s]:%d' "$1" "$2"
    else
        printf '%s:%d' "$1" "$2"
    fi
}

if ! (which jq &>/dev/null || install jq); then
    echo "JQ not installed and not installable.  The script jig requires it to function"
    exit 1
fi
`)

func putLog(uuid *strfmt.UUID, buf *bytes.Buffer) error {
	_, err := session.Jobs.PutJobLog(jobs.NewPutJobLogParams().WithUUID(*uuid).WithBody(buf), basicAuth)
	if err != nil {
		fmt.Printf("Failed to log to job log, %s: %v\n", uuid.String(), err)
	}
	return err
}

func Log(uuid *strfmt.UUID, echo bool, s string, args ...interface{}) error {
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, s, args...)
	if echo {
		fmt.Printf(s, args...)
	}
	return putLog(uuid, buf)
}

func writeStringToFile(filename, content string) error {
	dir := path.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	fo, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(content))
	if err != nil {
		return err
	}
	return nil
}

func markJob(uuid, state string, ops ModOps) error {
	j := fmt.Sprintf("{\"State\": \"%s\"}", state)
	if _, err := Update(uuid, j, ops, false); err != nil {
		fmt.Printf("Error marking job, %s, as %s: %v, continuing\n", uuid, state, err)
		return err
	}
	return nil
}

func markMachineRunnable(uuid string, ops ModOps) error {
	if _, err := Update(uuid, `{"Runnable": true}`, ops, false); err != nil {
		fmt.Printf("Error marking machine as runnable: %v, continuing to wait for runnable...\n", err)
		return err
	}
	return nil
}

type CommandRunner struct {
	name string
	uuid *strfmt.UUID

	cmd      *exec.Cmd
	stderr   io.ReadCloser
	stdout   io.ReadCloser
	stdin    io.WriteCloser
	task     *models.Task
	finished chan bool
}

func (cr *CommandRunner) ReadLog() {
	// read command's stderr line by line - for logging
	in := bufio.NewScanner(cr.stderr)
	for in.Scan() {
		putLog(cr.uuid, bytes.NewBuffer(in.Bytes()))
	}
	cr.finished <- true
}

func (cr *CommandRunner) ReadReply() {
	// read command's stdout line by line - for replies
	in := bufio.NewScanner(cr.stdout)
	for in.Scan() {
		putLog(cr.uuid, bytes.NewBuffer(in.Bytes()))
	}
	cr.finished <- true
}

func (cr *CommandRunner) Run() (failed, incomplete, reboot, poweroff bool) {
	// Start command running
	err := cr.cmd.Start()
	if err != nil {
		failed = true
		reboot = false
		Log(cr.uuid, true, "Command %s failed to start: %v\n", cr.name, err)
		return
	}

	// Wait for readers to exit
	<-cr.finished
	<-cr.finished

	err = cr.cmd.Wait()
	if exiterr, ok := err.(*exec.ExitError); ok {
		// The program has exited with an exit code != 0

		// This works on both Unix and Windows. Although package
		// syscall is generally platform dependent, WaitStatus is
		// defined for both Unix and Windows and in both cases has
		// an ExitStatus() method with the same signature.
		status, statusOK := exiterr.Sys().(syscall.WaitStatus)
		if !statusOK {
			if err != nil {
				failed = true
				reboot = false
			} else {
				failed = false
				reboot = false
			}
		} else {
			sane := false
			if cr.task.Meta != nil {
				flagStr, ok := cr.task.Meta["feature-flags"]
				if ok {
					for _, testFlag := range strings.Split(flagStr, ",") {
						if "sane-exit-codes" == strings.TrimSpace(testFlag) {
							sane = true
							break
						}
					}
				}
			}
			if !sane {
				if st, err := os.Stat(".sane-exit-codes"); err == nil && st.Mode().IsRegular() {
					sane = true
				}
			}
			code := uint(status.ExitStatus())
			if sane {
				// codes can be between 0 and 255
				// if the low bits are not 0, the command failed.
				failed = code&^224 > uint(0)
				// If the high bit is set, the command was incomplete and the
				// the current task pointer should not be advanced.
				incomplete = code&128 > uint(0)
				// If we need a reboot, set bit 6
				reboot = code&64 > uint(0)
				// If we need to poweroff, set bit 5.  Reboot wins if it is set.
				poweroff = code&32 > uint(0)
			} else {
				switch code {
				case 0:
				case 1:
					reboot = true
				case 2:
					incomplete = true
				case 3:
					incomplete = true
					reboot = true
				default:
					failed = true
				}
			}
		}
	} else {
		if err != nil {
			failed = true
			reboot = false
		} else {
			failed = false
			reboot = false
		}
	}
	Log(cr.uuid, true, "Command %s: failed: %v, incomplete: %v, reboot: %v, poweroff: %v\n",
		cr.name,
		failed,
		incomplete,
		reboot,
		poweroff)
	// Remove script
	os.Remove(cr.cmd.Path)

	return
}

func NewCommandRunner(uuid *strfmt.UUID, name, content, loc string) (*CommandRunner, error) {
	answer := &CommandRunner{name: name, uuid: uuid}
	if err := ioutil.WriteFile(path.Join(loc, "script"), []byte(content), 0700); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(path.Join(loc, "helper"), cmdHelper, 0400); err != nil {
		return nil, err
	}
	answer.cmd = exec.Command("./script")
	answer.cmd.Dir = loc
	answer.cmd.Env = append(os.Environ(), "RS_TASK_DIR="+loc, "RS_RUNNER_DIR="+runnerDir)
	var err2 error
	answer.stderr, err2 = answer.cmd.StderrPipe()
	if err2 != nil {
		return nil, err2
	}
	answer.stdout, err2 = answer.cmd.StdoutPipe()
	if err2 != nil {
		return nil, err2
	}
	answer.stdin, err2 = answer.cmd.StdinPipe()
	if err2 != nil {
		return nil, err2
	}

	answer.finished = make(chan bool, 2)
	go answer.ReadLog()
	go answer.ReadReply()

	return answer, nil
}

func runContent(uuid *strfmt.UUID, action *models.JobAction, task *models.Task) (failed, incomplete, reboot, poweroff bool) {
	tmpDir, err := ioutil.TempDir(runnerDir, *task.Name+"-")
	if err != nil {
		Log(uuid, true, "Could not create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	Log(uuid, false, "Starting Content Execution for: %s\n", *action.Name)

	runner, err := NewCommandRunner(uuid, *action.Name, *action.Content, tmpDir)
	if err != nil {
		failed = true
		Log(uuid, true, "Creating command %s failed: %v\n", *action.Name, err)
		return
	}
	runner.task = task
	return runner.Run()
}

func processJobsCommand() *cobra.Command {
	mo := &MachineOps{CommonOps{Name: "machines", SingularName: "machine"}}
	jo := &JobOps{CommonOps{Name: "jobs", SingularName: "job"}}
	so := &StageOps{CommonOps{Name: "stages", SingularName: "stage"}}
	to := &TaskOps{CommonOps{Name: "tasks", SingularName: "task"}}

	command := &cobra.Command{
		Use:   "processjobs [id]",
		Short: "For the given machine, process pending jobs until done.",
		Long: `
For the provided machine, identified by UUID, process the task list on
that machine until an error occurs or all jobs are complete.  Upon
completion, optionally wait for additional jobs as specified by
the stage runner wait flag.
`,
		RunE: func(c *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("%v requires at least 1 argument", c.UseLine())
			}
			if len(args) > 1 {
				return fmt.Errorf("%v requires at most 1 arguments", c.UseLine())
			}
			dumpUsage = false

			uuid := args[0]

			var machine *models.Machine
			if obj, err := Get(uuid, mo); err != nil {
				return generateError(err, "Error getting machine")
			} else {
				machine = obj.(*models.Machine)
			}
			var err error
			runnerDir, err = ioutil.TempDir("", "runner-")
			if err != nil {
				return generateError(err, "Error making temp runner dir")
			}
			defer os.RemoveAll(runnerDir)

			fmt.Printf("Processing jobs for %s\n", uuid)

			// Get Current Job and mark it failed if it is running.
			if obj, err := Get(machine.CurrentJob.String(), jo); err == nil {
				job := obj.(*models.Job)
				// If job is running or created, mark it as failed
				if *job.State == "running" || *job.State == "created" {
					markJob(machine.CurrentJob.String(), "failed", jo)
				}
			}

			// Mark Machine runnnable
			markMachineRunnable(machine.UUID.String(), mo)

			did_job := false
			for {
				// Wait for machine to be runnable.
				if answer, err := mo.DoWait(machine.UUID.String(), "Runnable", "true", 100000000); err != nil {
					fmt.Printf("Error waiting for machine to be runnable: %v, try again...\n", err)
					time.Sleep(5 * time.Second)
					continue
				} else if answer == "timeout" {
					fmt.Printf("Waiting for machine runnable returned with, %s, trying again.\n", answer)
					continue
				} else if answer == "interrupt" {
					fmt.Printf("User interrupted the wait, exiting ...\n")
					break
				}

				// Create a job for tasks
				var job *models.Job
				if obj, err := jo.Create(&models.Job{Machine: machine.UUID}); err != nil {
					fmt.Printf("Error creating a job for machine: %v, continuing\n", err)
					time.Sleep(5 * time.Second)
					continue
				} else {
					if obj == nil {
						if did_job {
							fmt.Println("Jobs finished")
							did_job = false
						}

						wait := false
						if obj, err := Get(uuid, mo); err == nil {
							machine = obj.(*models.Machine)

							if sobj, err := Get(machine.Stage, so); err == nil {
								stage := sobj.(*models.Stage)
								wait = stage.RunnerWait
							}
						}

						if wait {
							// Wait for new jobs - XXX: Web socket one day.
							// Create a not equal waiter
							time.Sleep(5 * time.Second)
							continue
						} else {
							break
						}
					}
					job = obj.(*models.Job)
				}
				did_job = true
				var task *models.Task

				if obj, err := Get(job.Task, to); err != nil {
					Log(job.UUID, true, "Error loading task content: %v, continuing", err)
					markJob(job.UUID.String(), "failed", jo)
				} else {
					task = obj.(*models.Task)
				}

				// Get the job data
				var list []*models.JobAction
				if resp, err := session.Jobs.GetJobActions(jobs.NewGetJobActionsParams().WithUUID(*job.UUID), basicAuth); err != nil {
					Log(job.UUID, true, "Error loading task content: %v, continuing", err)
					markJob(job.UUID.String(), "failed", jo)
					continue
				} else {
					list = resp.Payload
				}

				// Mark job as running
				if _, err := Update(job.UUID.String(), `{"State": "running"}`, jo, false); err != nil {
					Log(job.UUID, true, "Error marking job as running: %v, continue\n", err)
					markJob(job.UUID.String(), "failed", jo)
					continue
				}
				fmt.Printf("Starting Task: %s (%s)\n", job.Task, job.UUID.String())

				failed := false
				incomplete := false
				reboot := false
				poweroff := false
				state := "finished"

				for _, action := range list {
					event := &models.Event{Time: strfmt.DateTime(time.Now()), Type: "jobs", Action: "action_start", Key: job.UUID.String(), Object: fmt.Sprintf("Starting task: %s, template: %s", job.Task, *action.Name)}

					if _, err := session.Events.PostEvent(events.NewPostEventParams().WithBody(event), basicAuth); err != nil {
						fmt.Printf("Error posting event: %v\n", err)
					}

					// Excute task
					if *action.Path == "" {
						fmt.Printf("Running Task Template: %s\n", *action.Name)
						failed, incomplete, reboot, poweroff = runContent(job.UUID, action, task)
					} else {
						fmt.Printf("Putting Content in place for Task Template: %s\n", *action.Name)
						if err := writeStringToFile(*action.Path, *action.Content); err != nil {
							failed = true
							Log(job.UUID, true, "Task Template: %s - Copying contents to %s failed\n%v", *action.Name, *action.Path, err)
						} else {
							Log(job.UUID, true, "Task Template: %s - Copied contents to %s successfully\n", *action.Name, *action.Path)
						}
					}

					if failed {
						state = "failed"
					} else if incomplete {
						state = "incomplete"
					}

					fmt.Printf("Task Template , %s, %s\n", *action.Name, state)
					event = &models.Event{Time: strfmt.DateTime(time.Now()), Type: "jobs", Action: "action_stop", Key: job.UUID.String(), Object: fmt.Sprintf("Finished task: %s, template: %s, state: %s", job.Task, *action.Name, state)}

					if _, err := session.Events.PostEvent(events.NewPostEventParams().WithBody(event), basicAuth); err != nil {
						fmt.Printf("Error posting event: %v\n", err)
					}

					if failed || incomplete || reboot || poweroff {
						break
					}
				}

				fmt.Printf("Task: %s %s\n", job.Task, state)
				markJob(job.UUID.String(), state, jo)
				// Loop back and wait for the machine to get marked runnable again

				if reboot {
					if actuallyPowerThings {
						_, err := exec.Command("reboot").Output()
						if err != nil {
							Log(job.UUID, false, "Failed to issue reboot\n")
						}
					} else {
						Log(job.UUID, true, "Would have rebooted")
					}
				}
				if poweroff {
					if actuallyPowerThings {
						_, err := exec.Command("poweroff").Output()
						if err != nil {
							Log(job.UUID, false, "Failed to issue poweroff\n")
						}
					} else {
						Log(job.UUID, true, "Would have powered down")
					}
				}

				// If we failed, should we exit
				if exitOnFailure && failed {
					return fmt.Errorf("Task failed, exiting ...\n")
				}
			}

			return nil
		},
	}
	command.Flags().BoolVar(&exitOnFailure, "exit-on-failure", false, "Exit on failure of a task")

	return command
}
