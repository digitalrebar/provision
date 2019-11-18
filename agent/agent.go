package agent

import (
	"bufio"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/VictorLowther/jsonpatch2"

	"github.com/VictorLowther/jsonpatch2/utils"
	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/shirou/gopsutil/host"
)

type state int

const (
	AGENT_INIT = state(iota)
	AGENT_WAIT_FOR_RUNNABLE
	AGENT_RUN_TASK
	AGENT_WAIT_FOR_CHANGE_STAGE
	AGENT_CHANGE_STAGE
	AGENT_EXIT
	AGENT_REBOOT
	AGENT_POWEROFF
	AGENT_KEXEC
)

type si struct {
	BootTime uint64
	Machine  *models.Machine
}

// Agent implements a new machine agent structured as a finite
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
type Agent struct {
	state                                     state
	waitTimeout                               time.Duration
	client                                    *api.Client
	events                                    *api.EventStream
	machine                                   *models.Machine
	runnerDir, chrootDir, stateDir            string
	context                                   string
	doPower, exitOnNotRunnable, exitOnFailure bool
	logger                                    io.Writer
	bootTime                                  uint64
	err                                       error
	task                                      *runner
	taskMux                                   *sync.Mutex
	exitNow                                   bool
	kill                                      chan error
}

func (a *Agent) saveState() error {
	if a.stateDir == "" {
		return nil
	}
	saveFile := path.Join(a.stateDir, a.machine.Key()+".state.new")
	fi, err := os.Create(saveFile)
	if err != nil {
		return err
	}
	defer fi.Close()
	er := gzip.NewWriter(fi)
	defer er.Close()
	enc := json.NewEncoder(er)
	ss := si{
		BootTime: a.bootTime,
		Machine:  a.machine,
	}
	if err := enc.Encode(&ss); err != nil {
		return err
	}
	er.Flush()
	fi.Sync()
	return os.Rename(saveFile, strings.TrimSuffix(saveFile, ".new"))
}

func (a *Agent) tmpDir() string {
	return path.Dir(a.runnerDir)
}

// New creates a new FSM based Machine Agent that starts out in
// the AGENT_INIT state.
func New(c *api.Client, m *models.Machine,
	exitOnNotRunnable, exitOnFailure, actuallyPowerThings bool,
	logger io.Writer) (*Agent, error) {
	res := &Agent{
		state:             AGENT_INIT,
		client:            c,
		machine:           m,
		doPower:           actuallyPowerThings,
		exitOnFailure:     exitOnFailure,
		exitOnNotRunnable: exitOnNotRunnable,
		logger:            logger,
		waitTimeout:       1 * time.Hour,
		taskMux:           &sync.Mutex{},
	}
	if res.logger == nil {
		res.logger = os.Stderr
	}
	bt, err := host.BootTime()
	if err != nil {
		return nil, err
	}
	res.bootTime = bt
	return res, nil
}

// logf is a helper function to make logging of Agent actions a bit
// easier.
func (a *Agent) logf(f string, args ...interface{}) {
	fmt.Fprintf(a.logger, f, args...)
}

// Timeout allows you to change how long the Agent will wait for an
// event from dr-provision from the default of 1 hour.
func (a *Agent) Timeout(t time.Duration) *Agent {
	a.waitTimeout = t
	return a
}

func (a *Agent) StateLoc(s string) *Agent {
	a.stateDir = s
	return a
}

func (a *Agent) Context(s string) *Agent {
	a.context = s
	return a
}

func (a *Agent) RunnerDir(s string) *Agent {
	a.runnerDir = s
	return a
}

func (a *Agent) markNotRunnable() {
	if !(a.machine.Context == "" && a.context == "") {
		return
	}
	m := &models.Machine{}
	p := jsonpatch2.Patch{
		{Op: "test", Path: "/Context", Value: ""},
		{Op: "replace", Path: "/Runnable", Value: false},
	}
	if err := a.client.Req().Patch(p).UrlForM(a.machine).Do(m); err != nil {
		a.logf("Failed to mark machine not runnable: %v\n", err)
	}
}

func (a *Agent) power(cmdLine string) error {
	if !(a.doPower && a.context == "") {
		a.state = AGENT_EXIT
		return nil
	}
	if info, err := a.client.Info(); err == nil && !info.HasFeature("auto-boot-target") {
		var actionObj interface{}
		if err := a.client.Req().Get().
			UrlForM(a.machine, "actions", "nextbootpxe").Do(&actionObj); err == nil {
			emptyMap := map[string]interface{}{}
			a.client.Req().Post(emptyMap).
				UrlForM(a.machine, "actions", "nextbootpxe").Do(nil)
		}
	}
	a.markNotRunnable()
	cmd := exec.Command(cmdLine)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if cmd.Run() == nil {
		os.Exit(0)
	}
	return fmt.Errorf("Failed to %s", cmdLine)
}

func (a *Agent) exitOrSleep() {
	if a.exitOnFailure {
		a.state = AGENT_EXIT
	} else {
		time.Sleep(30 * time.Second)
	}
}

func (a *Agent) initOrExit() {
	if a.exitOnFailure {
		a.state = AGENT_EXIT
	} else {
		a.state = AGENT_INIT
		time.Sleep(5 * time.Second)
	}
}

// init resets the Machine Agent back to its initial state.  This
// consists of marking any current running jobs as Failed and
// reopening the event stream from dr-provision.
func (a *Agent) init() {
	a.taskMux.Lock()
	defer a.taskMux.Unlock()
	if a.err != nil {
		a.err = nil
	}
	if a.events != nil {
		a.events.Close()
		a.events = nil
	}
	var err error
	currentJob := &models.Job{Uuid: a.machine.CurrentJob}
	if a.client.Req().Fill(currentJob) == nil && currentJob.Context == a.context {
		// Only reset job state if we were responsible for creating it in the first place.
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
		a.logf("MachineAgent: error attaching to event stream: %v", err)
		a.exitOrSleep()
		return
	}
	a.state = AGENT_WAIT_FOR_RUNNABLE
}

func kexecLoaded() bool {
	buf, err := ioutil.ReadFile("/sys/kernel/kexec_loaded")
	return err == nil && string(buf)[0] == '1'
}

func (a *Agent) loadKexec() {
	kexecOk := false
	if err := a.client.Req().
		UrlFor("machines", a.machine.UUID(), "params", "kexec-ok").
		Params("aggregate", "true").
		Do(&kexecOk); err != nil {
		a.logf("kexec: Could not find kexec-ok\n")
		return
	}
	if !kexecOk {
		a.logf("kexec: kexec-ok is false\n")
		return
	}
	a.logf("Machine has kexec-ok param set\n")
	if runtime.GOOS != "linux" {
		a.logf("kexec: Not running on Linux\n")
		return
	}
	a.logf("Running on Linux\n")
	if _, err := exec.LookPath("kexec"); err != nil {
		a.logf("kexec: kexec command not installed\n")
		return
	}
	if kexecLoaded() {
		return
	}
	tmpDir, err := ioutil.TempDir(a.tmpDir(), "drp-agent-kexec")
	if err != nil {
		a.logf("Failed to make tmpdir for kexec\n")
		return
	}
	defer os.RemoveAll(tmpDir)
	kTmpl, err := a.client.File("machines", a.machine.UUID(), "kexec")
	if err != nil {
		a.logf("Failed to fetch kexec information: %v\n", err)
		return
	}
	a.logf("kexec info fetched\n")
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
			a.logf("Failed to fetch %s\n", parts[1])
			return
		}
		defer resp.Body.Close()
		outPath := path.Join(tmpDir, path.Base(parts[1]))
		out, err := os.Create(outPath)
		if err != nil {
			a.logf("Failed to create %s\n", outPath)
			return
		}
		if _, err := io.Copy(out, resp.Body); err != nil {
			a.logf("Failed to save %s\n", outPath)
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
		a.logf("No kernel found\n")
		return
	}
	if len(initrds) > 1 {
		a.logf("kexec: Too many initrds\n")
		return
	}
	kOpts, err := ioutil.ReadFile("/proc/cmdline")
	if err != nil {
		a.logf("kexec: No /proc/cmdline\n")
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
			a.logf("Could not determine nic we booted from")
			return
		}
	}
	a.logf("kernel:%s initrd:%s\n", kernel, initrds[0])
	a.logf("cmdline: %s\n", cmdline)
	cmd := exec.Command("/sbin/kexec", "-l", kernel, "--initrd="+initrds[0], "--command-line="+cmdline)
	if err := cmd.Run(); err != nil {
		return
	}
	a.logf("kexec info staged\n")
	return
}

func (a *Agent) doKexec() {
	if a.context == "" {
		a.state = AGENT_REBOOT
	} else {
		a.state = AGENT_EXIT
		return
	}
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

func (a *Agent) rebootOrExit(autoKexec bool) {
	a.markNotRunnable()
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
// conditions.  Once the conditions are met, the agent may transition
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
func (a *Agent) waitOn(m *models.Machine, cond api.TestFunc) {
	// No matter what else happens, we only respond when:
	// * The machine is available, and
	// * The ancilliary condition is met, and
	// * Either the machine context matches the one the agent cares about, or
	// * The bootenv changed.
	found, err := a.events.WaitFor(m,
		api.AndItems(
			api.EqualItem("Available", true),
			api.OrItems(
				api.EqualItem("Context", a.context),
				api.NotItem(api.EqualItem("BootEnv", a.machine.BootEnv))),
			cond),
		a.waitTimeout)
	if err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	a.logf("Wait: finished with %s\n", found)
	switch found {
	case "timeout":
		if a.exitOnNotRunnable {
			a.state = AGENT_EXIT
			return
		}
	case "interrupt":
		a.state = AGENT_EXIT
	case "complete":
		if m.BootEnv != a.machine.BootEnv && a.context == "" {
			a.rebootOrExit(true)
		} else if a.context == m.Context {
			if m.Runnable {
				a.state = AGENT_RUN_TASK
			} else {
				a.state = AGENT_WAIT_FOR_RUNNABLE
			}
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

// waitRunnable has waitOn wait for the Machine to become runnable.
func (a *Agent) waitRunnable() {
	m := models.Clone(a.machine).(*models.Machine)
	a.logf("Waiting on machine to become runnable\n")
	a.waitOn(m, api.EqualItem("Runnable", true))
}

// runTask attempts to run the next task on the Machine.  It may
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
func (a *Agent) runTask() {
	var err error
	a.taskMux.Lock()
	if a.exitNow {
		a.taskMux.Unlock()
		return
	}
	a.task, err = newRunner(a, a.machine, a.runnerDir, a.chrootDir, a.logger)
	a.taskMux.Unlock()
	if err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	if a.task == nil {
		if a.chrootDir != "" {
			a.logf("Current tasks finished, exiting chroot\n")
			a.state = AGENT_EXIT
			return
		}
		if a.machine.Workflow == "" {
			a.logf("Current tasks finished, check to see if stage needs to change\n")
			a.state = AGENT_CHANGE_STAGE
			return
		}
		a.logf("Current tasks finished, wait for stage or bootenv to change\n")
		a.state = AGENT_WAIT_FOR_CHANGE_STAGE
		return
	}
	defer func() { a.taskMux.Lock(); defer a.taskMux.Unlock(); a.task = nil }()
	a.logf("Runner created for task %s:%s:%s (%d:%d)\n",
		a.task.j.Workflow,
		a.task.j.Stage,
		a.task.j.Task,
		a.task.j.CurrentIndex,
		a.task.j.NextIndex)
	if a.task.wantChroot {
		a.chrootDir = a.task.jobDir
		a.state = AGENT_WAIT_FOR_RUNNABLE
		a.task.Close()
		return
	}
	if err := a.task.run(); err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	a.state = AGENT_WAIT_FOR_RUNNABLE
	if a.task.t != nil {
		defer a.task.Close()
		if a.task.reboot {
			a.task.log("Task signalled runner to reboot")
			a.rebootOrExit(false)
		} else if a.task.poweroff {
			a.task.log("Task signalled runner to poweroff")
			a.state = AGENT_POWEROFF
		} else if a.task.stop {
			a.task.log("Task signalled runner to stop")
			a.state = AGENT_EXIT
		} else if a.task.failed {
			a.task.log("Task signalled that it failed")
			if a.exitOnFailure {
				a.state = AGENT_EXIT
			}
		}
		if a.task.incomplete {
			a.task.log("Task signalled that it was incomplete")
		} else if !a.task.failed {
			a.task.log("Task signalled that it finished normally")
		}
	}
}

// waitChangeStage has waitOn wait for any of the following on the
// machine to change:
//
// * The CurrentTask index
// * The task list
// * The Runnable flag
// * The boot environment
// * The stage
func (a *Agent) waitChangeStage() {
	m := models.Clone(a.machine).(*models.Machine)
	a.logf("Waiting for system to be runnable and for stage or current tasks to change\n")
	a.waitOn(m,
		api.OrItems(api.NotItem(api.EqualItem("CurrentTask", m.CurrentTask)),
			api.NotItem(api.EqualItem("Tasks", m.Tasks)),
			api.NotItem(api.EqualItem("Runnable", m.Runnable)),
			api.NotItem(api.EqualItem("BootEnv", m.BootEnv)),
			api.NotItem(api.EqualItem("Stage", m.Stage))))
}

// changeStage handles determining what to do when the Agent runs out
// of tasks to run in the current Stage.  It may transition to the following states:
//
// * AGENT_WAIT_FOR_CHANGE_STAGE if there is no next stage for this
//   machine in the change stage map and it is not in an -install
//   bootenv.
//
// * AGENT_EXIT if there is no next stage in the change stage map and
//   the machine is in an -install bootenv.  In this case, changeStage
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
func (a *Agent) changeStage() {
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
	a.logf("Changing stage from %s to %s\n", a.machine.Stage, nextStage)
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

func (a *Agent) loadState() {
	if a.stateDir == "" {
		return
	}
	stateFile, err := os.Open(path.Join(a.stateDir, a.machine.Key()+".state"))
	if err != nil {
		return
	}
	defer stateFile.Close()
	var ss si
	gr, err := gzip.NewReader(stateFile)
	if err != nil {
		return
	}
	defer gr.Close()
	dec := json.NewDecoder(gr)
	if err = dec.Decode(&ss); err != nil {
		return
	}
	if ss.BootTime == a.bootTime && a.machine.Key() == ss.Machine.Key() {
		if a.machine.BootEnv != ss.Machine.BootEnv && a.context == "" {
			a.rebootOrExit(true)
		}
		a.machine = ss.Machine
		return
	}
}

func (a *Agent) Kill() error {
	a.logf("Agent signalled to exit")
	a.taskMux.Lock()
	a.kill = make(chan error)
	a.exitNow = true
	a.events.Kill()
	if a.task != nil {
		a.task.log("Agent signalled to exit")
		a.task.kill()
	}
	a.taskMux.Unlock()
	return <-a.kill
}

// Run kicks off the state machine for this agent.
func (a *Agent) Run() error {
	if a.context == "" && (a.machine.HasFeature("original-change-stage") ||
		!a.machine.HasFeature("change-stage-v2")) {
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
	if a.runnerDir == "" {
		a.runnerDir = os.Getenv("RS_RUNNER_DIR")
	}
	if a.runnerDir == "" {
		var td string
		if err := a.client.Req().UrlForM(a.machine, "params", "runner-tmpdir").Params("aggregate", "true").Do(&td); err == nil && td != "" {
			if err = mktd(td); err != nil {
				return err
			}
		}
		runnerDir, err := ioutil.TempDir(td, "runner-")
		if err != nil {
			return err
		}
		a.runnerDir = runnerDir
	}
	a.loadState()
	for {
		a.taskMux.Lock()
		if a.exitNow {
			a.state = AGENT_EXIT
		}
		a.taskMux.Unlock()
		if err := os.MkdirAll(a.runnerDir, 0755); err != nil {
			return err
		}
		if err := a.client.MakeProxy(path.Join(a.runnerDir, ".sock")); err != nil {
			return err
		}
		switch a.state {
		case AGENT_INIT:
			a.logf("Agent in init\n")
			a.init()
		case AGENT_WAIT_FOR_RUNNABLE:
			a.logf("Agent waiting for tasks\n")
			a.waitRunnable()
		case AGENT_RUN_TASK:
			a.logf("Agent running task\n")
			a.runTask()
		case AGENT_WAIT_FOR_CHANGE_STAGE:
			a.logf("Agent waiting stage change\n")
			a.waitChangeStage()
		case AGENT_CHANGE_STAGE:
			a.logf("Agent changing stage\n")
			a.changeStage()
		case AGENT_EXIT:
			if a.chrootDir != "" {
				a.logf("Agent exiting chroot %s\n", a.chrootDir)
				a.chrootDir = ""
				a.waitRunnable()
			} else {
				a.logf("Agent exiting\n")
				a.taskMux.Lock()
				if a.exitNow {
					a.kill <- a.err
				}
				a.taskMux.Unlock()
				return a.err
			}
		case AGENT_KEXEC:
			a.logf("Attempting to kexec\n")
			a.doKexec()
		case AGENT_REBOOT:
			a.logf("Agent rebooting\n")
			return a.power("reboot")
		case AGENT_POWEROFF:
			a.logf("Agent powering off\n")
			return a.power("poweroff")
		default:
			a.logf("Unknown agent state %d\n", a.state)
			panic("unreachable")
		}
		if a.err != nil {
			a.logf("Error during run: %v\n", a.err)
		}
		if err := a.saveState(); err != nil {
			a.logf("Error saving state: %v", err)
		}
	}
}
