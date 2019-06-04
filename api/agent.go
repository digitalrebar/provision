package api

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/models"
)

type AgentState int

const (
	AGENT_INIT = AgentState(iota)
	AGENT_WAIT_FOR_RUNNABLE
	AGENT_RUN_TASK
	AGENT_WAIT_FOR_CHANGE_STAGE
	AGENT_CHANGE_STAGE
	AGENT_EXIT
	AGENT_REBOOT
	AGENT_POWEROFF
	AGENT_KEXEC
)

// MachineAgent implements a new machine agent structured as a finite
// state machine.  There is one important behavioural change to the
// behaviour of the runner that may impact how workflows are built:
//
// The RunnerWait flag in stages is no longer honored.  Instead, the
// agent will wait by default, unless overridden by the following
// conditions, in order of priority:
//
// * The next stage has the Reboot flag set.
//
// * The change-stage/map entry for the next stage has a Stop, Reboot,
//   or Poweroff clause.
//
// * The machine is currently in a bootenv that ends in -install and
//   there is nothing else to do, in which case the runner will exit
//
//
// Additionally, this agent will automatically reboot the system when
// it detects that the machine's boot environment has changed, unless
// the machine is in an OS install, in which case the agent will exit.
type MachineAgent struct {
	state                                     AgentState
	waitTimeout                               time.Duration
	client                                    *Client
	events                                    *EventStream
	machine                                   *models.Machine
	runnerDir, chrootDir                      string
	doPower, exitOnNotRunnable, exitOnFailure bool
	logger                                    io.Writer
	err                                       error
}

// NewAgent creates a new FSM based Machine Agent that starts out in
// the AGENT_INIT state.
func (c *Client) NewAgent(m *models.Machine,
	exitOnNotRunnable, exitOnFailure, actuallyPowerThings bool,
	logger io.Writer) (*MachineAgent, error) {
	res := &MachineAgent{
		state:             AGENT_INIT,
		client:            c,
		machine:           m,
		doPower:           actuallyPowerThings,
		exitOnFailure:     exitOnFailure,
		exitOnNotRunnable: exitOnNotRunnable,
		logger:            logger,
		waitTimeout:       1 * time.Hour,
	}
	if res.logger == nil {
		res.logger = os.Stderr
	}
	rdExists := false
	res.runnerDir = os.Getenv("RS_RUNNER_DIR")
	if res.runnerDir != "" {
		if fi, err := os.Stat(res.runnerDir); err == nil {
			if fi.IsDir() {
				rdExists = true
			}
		}
	}
	if res.runnerDir == "" {
		var td string
		if err := c.Req().UrlForM(m, "params", "runner-tmpdir").Params("aggregate", "true").Do(&td); err == nil && td != "" {
			if err = mktd(td); err != nil {
				return nil, err
			}
			if err = os.Setenv("TMPDIR", td); err != nil {
				return nil, err
			}
			if err = os.Setenv("TMP", td); err != nil {
				return nil, err
			}
		}
		runnerDir, err := ioutil.TempDir("", "runner-")
		if err != nil {
			return nil, err
		}
		if err := os.Setenv("RS_RUNNER_DIR", runnerDir); err != nil {
			return nil, err
		}
		res.runnerDir = runnerDir
	}
	if !rdExists {
		if err := os.MkdirAll(res.runnerDir, 0755); err != nil {
			return nil, err
		}
	}
	return res, nil
}

// Logf is a helper function to make logging of Agent actions a bit
// easier.
func (a *MachineAgent) Logf(f string, args ...interface{}) {
	fmt.Fprintf(a.logger, f, args...)
}

// Timeout allows you to change how long the Agent will wait for an
// event from dr-provision from the default of 1 hour.
func (a *MachineAgent) Timeout(t time.Duration) *MachineAgent {
	a.waitTimeout = t
	return a
}

func (a *MachineAgent) power(cmdLine string) error {
	if !a.doPower {
		return nil
	}
	if !a.client.info.HasFeature("auto-boot-target") {
		var actionObj interface{}
		if err := a.client.Req().Get().
			UrlForM(a.machine, "actions", "nextbootpxe").Do(&actionObj); err == nil {
			emptyMap := map[string]interface{}{}
			a.client.Req().Post(emptyMap).
				UrlForM(a.machine, "actions", "nextbootpxe").Do(nil)
		}
	}
	cmd := exec.Command(cmdLine)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if cmd.Run() == nil {
		os.Exit(0)
	}
	return fmt.Errorf("Failed to %s", cmdLine)
}

func (a *MachineAgent) exitOrSleep() {
	if a.exitOnFailure {
		a.state = AGENT_EXIT
	} else {
		time.Sleep(30 * time.Second)
	}
}

func (a *MachineAgent) initOrExit() {
	if a.exitOnFailure {
		a.state = AGENT_EXIT
	} else {
		a.state = AGENT_INIT
		time.Sleep(5 * time.Second)
	}
}

// Init resets the Machine Agent back to its initial state.  This
// consists of marking any current running jobs as Failed and
// repoening the event stream from dr-provision.
func (a *MachineAgent) Init() {
	if a.err != nil {
		a.err = nil
	}
	if a.events != nil {
		a.events.Close()
		a.events = nil
	}
	var err error
	currentJob := &models.Job{Uuid: a.machine.CurrentJob}
	if a.client.Req().Fill(currentJob) == nil {
		if currentJob.State == "running" || currentJob.State == "created" {
			cj := models.Clone(currentJob).(*models.Job)
			cj.State = "failed"
			if _, a.err = a.client.PatchTo(currentJob, cj); a.err != nil {
				a.exitOrSleep()
				return
			}
		}
	}
	a.events, a.err = a.client.Events()
	if a.err != nil {
		a.Logf("MachineAgent: error attaching to event stream: %v", err)
		a.exitOrSleep()
		return
	}
	a.state = AGENT_WAIT_FOR_RUNNABLE
}

func kexecLoaded() bool {
	buf, err := ioutil.ReadFile("/sys/kernel/kexec_loaded")
	return err == nil && string(buf)[0] == '1'
}

func (a *MachineAgent) loadKexec() {
	kexecOk := false
	if err := a.client.Req().
		UrlFor("machines", a.machine.UUID(), "params", "kexec-ok").
		Params("aggregate", "true").
		Do(&kexecOk); err != nil {
		a.Logf("kexec: Could not find kexec-ok\n")
		return
	}
	if !kexecOk {
		a.Logf("kexec: kexec-ok is false\n")
		return
	}
	a.Logf("Machine has kexec-ok param set\n")
	if runtime.GOOS != "linux" {
		a.Logf("kexec: Not running on Linux\n")
		return
	}
	a.Logf("Running on Linux\n")
	if _, err := exec.LookPath("kexec"); err != nil {
		a.Logf("kexec: kexec command not installed\n")
		return
	}
	if kexecLoaded() {
		return
	}
	tmpDir, err := ioutil.TempDir("", "drp-agent-kexec")
	if err != nil {
		a.Logf("Failed to make tmpdir for kexec\n")
		return
	}
	defer os.RemoveAll(tmpDir)
	kTmpl, err := a.client.File("machines", a.machine.UUID(), "kexec")
	if err != nil {
		a.Logf("Failed to fetch kexec information: %v\n", err)
		return
	}
	a.Logf("kexec info fetched\n")
	defer kTmpl.Close()
	sc := bufio.NewScanner(kTmpl)
	var kernel, cmdline string
	initrds := []string{}
	for sc.Scan() {
		parts := strings.SplitN(sc.Text(), " ", 2)
		var resp *http.Response
		switch parts[0] {
		case "kernel", "initrd":
			resp, err = http.Get(parts[1])
		case "params":
			cmdline = parts[1]
			continue
		default:
			continue
		}
		if err != nil {
			a.Logf("Failed to fetch %s\n", parts[1])
			return
		}
		defer resp.Body.Close()
		outPath := path.Join(tmpDir, path.Base(parts[1]))
		out, err := os.Create(outPath)
		if err != nil {
			a.Logf("Failed to create %s\n", outPath)
			return
		}
		if _, err := io.Copy(out, resp.Body); err != nil {
			a.Logf("Failed to save %s\n", outPath)
			return
		}
		out.Sync()
		out.Close()
		switch parts[0] {
		case "kernel":
			kernel = outPath
		case "initrd":
			initrds = append(initrds, outPath)
		}
	}
	if kernel == "" {
		a.Logf("No kernel found\n")
		return
	}
	if len(initrds) > 1 {
		a.Logf("kexec: Too many initrds\n")
		return
	}
	kOpts, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		a.Logf("kexec: No /proc/cmdline\n")
		return
	}
	kexecOk = false
	for _, part := range strings.Split(string(kOpts), " ") {
		if strings.HasPrefix(part, "BOOTIF=") {
			kexecOk = true
			cmdline = cmdline + " " + part
			break
		}
	}
	if !kexecOk {
		v, vok := a.machine.Params["last-boot-macaddr"]
		macaddr, aok := v.(string)
		if aok && vok {
			cmdline = cmdline + " BOOTIF=" + macaddr
		} else {
			a.Logf("Could not determine nic we booted from")
			return
		}
	}
	a.Logf("kernel:%s initrd:%s\n", kernel, initrds[0])
	a.Logf("cmdline: %s\n", cmdline)
	cmd := exec.Command("/sbin/kexec", "-l", kernel, "--initrd="+initrds[0], "--command-line="+cmdline)
	if err := cmd.Run(); err != nil {
		return
	}
	a.Logf("kexec info staged\n")
	return
}

func (a *MachineAgent) doKexec() {
	a.state = AGENT_REBOOT
	var cmdErr error
	if _, err := exec.LookPath("systemctl"); err == nil {
		cmdErr = exec.Command("systemctl", "kexec").Run()
	} else if _, err = exec.LookPath("/etc/init.d/kexec"); err == nil {
		cmdErr = exec.Command("telinit", "6").Run()
	} else if err = exec.Command("grep", "-q", "kexec", "/etc/init.d/halt").Run(); err == nil {
		cmdErr = exec.Command("telinit", "6").Run()
	} else {
		cmdErr = exec.Command("kexec", "-e").Run()
	}
	if cmdErr == nil {
		time.Sleep(5 * time.Minute)
	}
}

func (a *MachineAgent) rebootOrExit(autoKexec bool) {
	if strings.HasSuffix(a.machine.BootEnv, "-install") {
		a.state = AGENT_EXIT
		return
	}
	if autoKexec {
		a.loadKexec()
	}
	if kexecLoaded() {
		a.state = AGENT_KEXEC
	} else {
		a.state = AGENT_REBOOT
	}
}

// waitOn waits for the machine to match the passed wait
// conitions.  Once the conditions are met, the agent may transition
// to the following states (in order of priority):
//
// * AGENT_EXIT if the machine wants to change from an -install
//   bootenv to a different bootenv
//
// * AGENT_REBOOT if the machine wants to change bootenvs
//
// * AGENT_RUN_TASK if the machine is runnable.
//
// * AGENT_WAIT_FOR_RUNNABLE if the machine is not runnable.
func (a *MachineAgent) waitOn(m *models.Machine, cond TestFunc) {
	found, err := a.events.WaitFor(m, AndItems(EqualItem("Available", true), cond), a.waitTimeout)
	if err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	a.Logf("Wait: finished with %s\n", found)
	switch found {
	case "timeout":
		if a.exitOnNotRunnable {
			a.state = AGENT_EXIT
			return
		}
	case "interrupt":
		a.state = AGENT_EXIT
	case "complete":
		if m.BootEnv != a.machine.BootEnv {
			a.rebootOrExit(true)
		} else if m.Runnable {
			a.state = AGENT_RUN_TASK
		} else {
			a.state = AGENT_WAIT_FOR_RUNNABLE
		}
	default:
		err := &models.Error{
			Type:  "AGENT_WAIT",
			Model: m.Prefix(),
			Key:   m.Key(),
		}
		err.Errorf("Unexpected return from WaitFor: %s", found)
		a.err = err
		a.initOrExit()
	}
	a.machine = m
}

// WaitRunnable has waitOn wait for the Machine to become runnable.
func (a *MachineAgent) WaitRunnable() {
	m := models.Clone(a.machine).(*models.Machine)
	a.Logf("Waiting on machine to become runnable\n")
	a.waitOn(m, EqualItem("Runnable", true))
}

// RunTask attempts to run the next task on the Machine.  It may
// transition to the following states:
//
// * AGENT_CHANGE_STAGE if there are no tasks to run.
//
// * AGENT_REBOOT if the task signalled that the machine should reboot
//
// * AGENT_POWEROFF if the task signalled that the machine should shut down
//
// * AGENT_EXIT if the task signalled that the agent should stop.
//
// * AGENT_WAIT_FOR_RUNNABLE if no other conditions were met.
func (a *MachineAgent) RunTask() {
	runner, err := NewTaskRunner(a.client, a.machine, a.runnerDir, a.chrootDir, a.logger)
	if err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	if runner == nil {
		if a.chrootDir != "" {
			a.Logf("Current tasks finished, exiting chroot\n")
			a.state = AGENT_EXIT
			return
		}
		if a.machine.Workflow == "" {
			a.Logf("Current tasks finished, check to see if stage needs to change\n")
			a.state = AGENT_CHANGE_STAGE
			return
		}
		a.Logf("Current tasks finished, wait for stage or bootenv to change\n")
		a.state = AGENT_WAIT_FOR_CHANGE_STAGE
		return
	}
	a.Logf("Runner created for task %s:%s:%s (%d:%d)\n",
		runner.j.Workflow,
		runner.j.Stage,
		runner.j.Task,
		runner.j.CurrentIndex,
		runner.j.NextIndex)
	if runner.wantChroot {
		a.chrootDir = runner.jobDir
		a.state = AGENT_WAIT_FOR_RUNNABLE
		runner.Close()
		return
	}
	if err := runner.Run(); err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	a.state = AGENT_WAIT_FOR_RUNNABLE
	if runner.t != nil {
		defer runner.Close()
		if runner.reboot {
			runner.Log("Task signalled runner to reboot")
			a.rebootOrExit(false)
		} else if runner.poweroff {
			runner.Log("Task signalled runner to poweroff")
			a.state = AGENT_POWEROFF
		} else if runner.stop {
			runner.Log("Task signalled runner to stop")
			a.state = AGENT_EXIT
		} else if runner.failed {
			runner.Log("Task signalled that it failed")
			if a.exitOnFailure {
				a.state = AGENT_EXIT
			}
		}
		if runner.incomplete {
			runner.Log("Task signalled that it was incomplete")
		} else if !runner.failed {
			runner.Log("Task signalled that it finished normally")
		}
	}
}

// WaitChangeStage has waitOn wait for any of the following on the
// machine to change:
//
// * The CurrentTask index
// * The task list
// * The Runnable flag
// * The boot environment
// * The stage
func (a *MachineAgent) WaitChangeStage() {
	m := models.Clone(a.machine).(*models.Machine)
	a.Logf("Waiting for system to be runnable and for stage or current tasks to change\n")
	a.waitOn(m,
		OrItems(NotItem(EqualItem("CurrentTask", m.CurrentTask)),
			NotItem(EqualItem("Tasks", m.Tasks)),
			NotItem(EqualItem("Runnable", m.Runnable)),
			NotItem(EqualItem("BootEnv", m.BootEnv)),
			NotItem(EqualItem("Stage", m.Stage))))
}

// ChangeStage handles determining what to do when the Agent runs out
// of tasks to run in the current Stage.  It may transition to the following states:
//
// * AGENT_WAIT_FOR_CHANGE_STAGE if there is no next stage for this
//   machine in the change stage map and it is not in an -install
//   bootenv.
//
// * AGENT_EXIT if there is no next stage in the change stage map and
//   the machine is in an -install bootenv.  In this case, ChangeStage
//   will set the machine stage to `local`.
//
// * AGENT_REBOOT if the next stage has the Reboot flag, the change
//   stage map has a Reboot specifier, or the next stage has a different
//   bootenv than the machine and the machine is not in an -install
//   bootenv
//
// * AGENT_EXIT if the machine is in an -install bootenv and the next
//   stage requires a different bootenv.
//
// * AGENT_POWEROFF if the change stage map wants to power the system
//   off after changing the stage.
//
// * AGENT_WAIT_FOR_RUNNABLE if no other condition applies
func (a *MachineAgent) ChangeStage() {
	var cmObj interface{}
	a.state = AGENT_WAIT_FOR_CHANGE_STAGE
	inInstall := strings.HasSuffix(a.machine.BootEnv, "-install")
	csMap := map[string]string{}
	csErr := a.client.Req().Get().
		UrlForM(a.machine, "params", "change-stage/map").
		Params("aggregate", "true").Do(&cmObj)
	if csErr == nil {
		if err := utils.Remarshal(cmObj, &csMap); err != nil {
			a.err = err
			a.initOrExit()
			return
		}
	}
	var nextStage, targetState string
	if ns, ok := csMap[a.machine.Stage]; ok {
		pieces := strings.SplitN(ns, ":", 2)
		nextStage = pieces[0]
		if len(pieces) == 2 {
			targetState = pieces[1]
		}
	}
	if nextStage == "" {
		if inInstall {
			nextStage = "local"
		} else {
			nextStage = a.machine.Stage
		}
	}
	if nextStage == a.machine.Stage {
		return
	}
	a.Logf("Changing stage from %s to %s\n", a.machine.Stage, nextStage)
	newStage := &models.Stage{}
	if err := a.client.FillModel(newStage, nextStage); err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	// Default behaviour for what to do for the next state
	if newStage.BootEnv == "" || newStage.BootEnv == a.machine.BootEnv {
		// If the bootenv has not changed, the machine will get a new task list.
		// Wait for the machine to be runnable if needed, and start running it.
		a.state = AGENT_WAIT_FOR_RUNNABLE
	} else {
		// The new stage wants a new bootenv.  Reboot into it to continue
		// processing.
		a.rebootOrExit(true)
	}
	if targetState != "" {
		// The change stage map is overriding our default behaviour.
		switch targetState {
		case "Reboot":
			a.rebootOrExit(true)
		case "Stop":
			a.state = AGENT_EXIT
		case "Shutdown":
			a.state = AGENT_POWEROFF
		}
	}
	if newStage.Reboot {
		// A reboot flag on the next stage forces an unconditional reboot.
		a.rebootOrExit(true)
	}
	newM := models.Clone(a.machine).(*models.Machine)
	newM.Stage = nextStage
	if _, err := a.client.PatchTo(a.machine, newM); err != nil {
		a.err = err
		a.initOrExit()
	}
}

// Run kicks off the state machine for this agent.
func (a *MachineAgent) Run() error {
	if a.machine.HasFeature("original-change-stage") ||
		!a.machine.HasFeature("change-stage-v2") {
		newM := models.Clone(a.machine).(*models.Machine)
		newM.Runnable = true
		if err := a.client.Req().PatchTo(a.machine, newM).Do(&newM); err == nil {
			a.machine = newM
		} else {
			res := &models.Error{
				Type:  "AGENT_WAIT",
				Model: a.machine.Prefix(),
				Key:   a.machine.Key(),
			}
			res.Errorf("Failed to mark machine runnable.")
			res.AddError(err)
			return res
		}
	}
	for {
		switch a.state {
		case AGENT_INIT:
			a.Logf("Agent in init\n")
			a.Init()
		case AGENT_WAIT_FOR_RUNNABLE:
			a.Logf("Agent waiting for tasks\n")
			a.WaitRunnable()
		case AGENT_RUN_TASK:
			a.Logf("Agent running task\n")
			a.RunTask()
		case AGENT_WAIT_FOR_CHANGE_STAGE:
			a.Logf("Agent waiting stage change\n")
			a.WaitChangeStage()
		case AGENT_CHANGE_STAGE:
			a.Logf("Agent changing stage\n")
			a.ChangeStage()
		case AGENT_EXIT:
			if a.chrootDir != "" {
				a.Logf("Agent exiting chroot %s\n", a.chrootDir)
				a.chrootDir = ""
				a.WaitRunnable()
			} else {
				a.Logf("Agent exiting\n")
				return a.err
			}
		case AGENT_KEXEC:
			a.Logf("Attempting to kexec\n")
			a.doKexec()
		case AGENT_REBOOT:
			a.Logf("Agent rebooting\n")
			return a.power("reboot")
		case AGENT_POWEROFF:
			a.Logf("Agent powering off\n")
			return a.power("poweroff")
		default:
			a.Logf("Unknown agent state %d\n", a.state)
			panic("unreachable")
		}
		if a.err != nil {
			a.Logf("Error during run: %v\n", a.err)
		}
	}
}
