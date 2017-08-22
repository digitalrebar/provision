package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
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

func Log(uuid *strfmt.UUID, s string) error {
	buf := bytes.NewBufferString(s)
	_, err := session.Jobs.PutJobLog(jobs.NewPutJobLogParams().WithUUID(*uuid).WithBody(buf), basicAuth)
	if err != nil {
		fmt.Printf("Failed to log to job log, %s: %v\n", uuid.String(), err)
	}
	return err
}

func writeStringToFile(path, content string) error {
	fo, err := os.Create(path)
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
	finished chan bool
}

func (cr *CommandRunner) ReadLog() {
	// read command's stderr line by line - for logging
	in := bufio.NewScanner(cr.stderr)
	for in.Scan() {
		Log(cr.uuid, in.Text())
	}
	cr.finished <- true
}

func (cr *CommandRunner) ReadReply() {
	// read command's stdout line by line - for replies
	in := bufio.NewScanner(cr.stdout)
	for in.Scan() {
		Log(cr.uuid, in.Text())
	}
	cr.finished <- true
}

func (cr *CommandRunner) Run() (failed, incomplete, reboot bool) {
	// Start command running
	err := cr.cmd.Start()
	if err != nil {
		failed = true
		reboot = false
		s := fmt.Sprintf("Command %s failed to start: %v\n", cr.name, err)
		fmt.Printf(s)
		Log(cr.uuid, s)
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
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			code := status.ExitStatus()
			switch code {
			case 0:
				failed = false
				reboot = false
				s := fmt.Sprintf("Command %s succeeded\n", cr.name)
				fmt.Printf(s)
				Log(cr.uuid, s)
			case 1:
				failed = false
				reboot = true
				s := fmt.Sprintf("Command %s succeeded (wants reboot)\n", cr.name)
				fmt.Printf(s)
				Log(cr.uuid, s)
			case 2:
				incomplete = true
				reboot = false
				s := fmt.Sprintf("Command %s incomplete\n", cr.name)
				fmt.Printf(s)
				Log(cr.uuid, s)
			case 3:
				incomplete = true
				reboot = true
				s := fmt.Sprintf("Command %s incomplete (wants reboot)\n", cr.name)
				fmt.Printf(s)
				Log(cr.uuid, s)
			default:
				failed = true
				reboot = false
				s := fmt.Sprintf("Command %s failed\n", cr.name)
				fmt.Printf(s)
				Log(cr.uuid, s)
			}
		}
	} else {
		if err != nil {
			failed = true
			reboot = false
			s := fmt.Sprintf("Command %s failed: %v\n", cr.name, err)
			fmt.Printf(s)
			Log(cr.uuid, s)
		} else {
			failed = false
			reboot = false
			s := fmt.Sprintf("Command %s succeeded\n", cr.name)
			fmt.Printf(s)
			Log(cr.uuid, s)
		}
	}

	// Remove script
	os.Remove(cr.cmd.Path)

	return
}

func NewCommandRunner(uuid *strfmt.UUID, name, content string) (*CommandRunner, error) {
	answer := &CommandRunner{name: name, uuid: uuid}

	// Make script file
	tmpFile, err := ioutil.TempFile(".", "script")
	if err != nil {
		return nil, err
	}
	if _, err := tmpFile.Write([]byte(content)); err != nil {
		return nil, err
	}
	path := "./" + tmpFile.Name()
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}
	os.Chmod(path, 0700)

	answer.cmd = exec.Command(path)

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

func runContent(uuid *strfmt.UUID, action *models.JobAction) (failed, incomplete, reboot bool) {

	Log(uuid, fmt.Sprintf("Starting Content Execution for: %s\n", *action.Name))

	runner, err := NewCommandRunner(uuid, *action.Name, *action.Content)
	if err != nil {
		failed = true
		s := fmt.Sprintf("Creating command %s failed: %v\n", *action.Name, err)
		fmt.Printf(s)
		Log(uuid, s)
	} else {
		failed, incomplete, reboot = runner.Run()
	}
	return
}

func processJobsCommand() *cobra.Command {
	mo := &MachineOps{CommonOps{Name: "machines", SingularName: "machine"}}
	jo := &JobOps{CommonOps{Name: "jobs", SingularName: "job"}}

	command := &cobra.Command{
		Use:   "processjobs [id] [wait]",
		Short: "For the given machine, process pending jobs until done.",
		Long: `
For the provided machine, identified by UUID, process the task list on
that machine until an error occurs or all jobs are complete.  Upon 
completion, optionally wait for additional jobs as specified by
the boolean wait flag.
`,
		RunE: func(c *cobra.Command, args []string) error {
			var err error
			if len(args) < 1 {
				return fmt.Errorf("%v requires at least 1 argument", c.UseLine())

			}
			if len(args) > 2 {
				return fmt.Errorf("%v requires at most 2 arguments", c.UseLine())
			}
			dumpUsage = false

			uuid := args[0]
			wait := false
			if len(args) == 2 {
				wait, err = strconv.ParseBool(args[1])
				if err != nil {
					return fmt.Errorf("Error reading wait argument: %v", err)
				}
			}

			waitStr := "will not wait for new jobs"
			if wait {
				waitStr = "will wait for new jobs"
			}

			var machine *models.Machine
			if obj, err := Get(uuid, mo); err != nil {
				return generateError(err, "Error getting machine")
			} else {
				machine = obj.(*models.Machine)
			}

			fmt.Printf("Processing jobs for %s (%s)\n", uuid, waitStr)

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

				// Get the job data
				var list []*models.JobAction
				if resp, err := session.Jobs.GetJobActions(jobs.NewGetJobActionsParams().WithUUID(*job.UUID), basicAuth); err != nil {
					fmt.Printf("Error loading task content: %v, continuing", err)
					markJob(job.UUID.String(), "failed", jo)
					continue
				} else {
					list = resp.Payload
				}

				// Mark job as running
				if _, err := Update(job.UUID.String(), `{"State": "running"}`, jo, false); err != nil {
					fmt.Printf("Error marking job as running: %v, continue\n", err)
					markJob(job.UUID.String(), "failed", jo)
					continue
				}
				fmt.Printf("Starting Task: %s (%s)\n", job.Task, job.UUID.String())

				failed := false
				incomplete := false
				reboot := false
				state := "finished"

				for _, action := range list {
					event := &models.Event{Time: strfmt.DateTime(time.Now()), Type: "jobs", Action: "action_start", Key: job.UUID.String(), Object: fmt.Sprintf("Starting task: %s, template: %s", job.Task, *action.Name)}

					if _, err := session.Events.PostEvent(events.NewPostEventParams().WithBody(event), basicAuth); err != nil {
						fmt.Printf("Error posting event: %v\n", err)
					}

					// Excute task
					if *action.Path == "" {
						fmt.Printf("Running Task Template: %s\n", *action.Name)
						failed, incomplete, reboot = runContent(job.UUID, action)
					} else {
						fmt.Printf("Putting Content in place for Task Template: %s\n", *action.Name)
						var s string
						if err := writeStringToFile(*action.Path, *action.Content); err != nil {
							failed = true
							s = fmt.Sprintf("Task Template: %s - Copying contents to %s failed\n%v", *action.Name, *action.Path, err)
						} else {
							s = fmt.Sprintf("Task Template: %s - Copied contents to %s successfully\n", *action.Name, *action.Path)
						}
						fmt.Printf(s)
						Log(job.UUID, s)
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

					if failed || incomplete || reboot {
						break
					}
				}

				fmt.Printf("Task: %s %s\n", job.Task, state)
				markJob(job.UUID.String(), state, jo)
				// Loop back and wait for the machine to get marked runnable again

				if reboot {
					_, err := exec.Command("reboot").Output()
					if err != nil {
						Log(job.UUID, "Failed to issue reboot\n")
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
