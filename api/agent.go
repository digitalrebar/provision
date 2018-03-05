package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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
)

type MachineAgent struct {
	state                                     AgentState
	client                                    *Client
	events                                    *EventStream
	machine                                   *models.Machine
	runnerDir                                 string
	doPower, exitOnNotRunnable, exitOnFailure bool
	logger                                    io.Writer
	err                                       error
}

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
	}
	if res.logger == nil {
		res.logger = os.Stderr
	}
	runnerDir, err := ioutil.TempDir("", "runner-")
	if err != nil {
		return nil, err
	}
	res.runnerDir = runnerDir
	return res, nil
}

func (a *MachineAgent) power(cmdLine string) error {
	if !a.doPower {
		return nil
	}
	var actionObj interface{}
	if err := a.client.Req().Get().
		UrlForM(a.machine, "actions", "nextbootpxe").Do(&actionObj); err == nil {
		emptyMap := map[string]interface{}{}
		a.client.Req().Post(emptyMap).
			UrlForM(a.machine, "actions", "nextbootpxe").Do(nil)
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

func (a *MachineAgent) Init() {
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
		log.Printf("MachineAgent: error attaching to event stream: %v", err)
		a.exitOrSleep()
		return
	}
	a.state = AGENT_WAIT_FOR_RUNNABLE
}

func (a *MachineAgent) waitOn(m *models.Machine, cond TestFunc) {
	found, err := a.events.WaitFor(m, cond, 1*time.Hour)
	if err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	switch found {
	case "timeout":
	case "interrupt":
		a.state = AGENT_EXIT
	case "complete":
		if m.BootEnv == a.machine.BootEnv {
			a.state = AGENT_RUN_TASK
		} else if strings.HasSuffix(m.BootEnv, "-install") {
			a.state = AGENT_EXIT
		} else {
			a.state = AGENT_REBOOT
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

func (a *MachineAgent) WaitRunnable() {
	m := models.Clone(a.machine).(*models.Machine)
	a.waitOn(m, EqualItem("Runnable", m.Runnable))
}

func (a *MachineAgent) RunTask() {
	runner, err := NewTaskRunner(a.client, a.machine, a.runnerDir, a.logger)
	if err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	if runner == nil {
		a.state = AGENT_CHANGE_STAGE
		return
	}
	if err := runner.Run(); err != nil {
		a.err = err
		a.initOrExit()
		return
	}
	a.state = AGENT_WAIT_FOR_RUNNABLE
	defer runner.Close()
	if runner.reboot {
		runner.Log("Task signalled runner to reboot")
		a.state = AGENT_REBOOT
	} else if runner.poweroff {
		runner.Log("Task signalled runner to poweroff")
		a.state = AGENT_POWEROFF
	} else if runner.stop {
		runner.Log("Task signalled runner to stop")
		a.state = AGENT_EXIT
	} else if runner.failed {
		runner.Log("Task signalled that it failed")
	}
	if runner.incomplete {
		runner.Log("Task signalled that it was incomplete")
	} else {
		runner.Log("Task signalled that it finished normally")
	}
}

func (a *MachineAgent) WaitChangeStage() {
	m := models.Clone(a.machine).(*models.Machine)
	a.waitOn(m,
		AndItems(EqualItem("Runnable", m.Runnable),
			OrItems(NotItem(EqualItem("CurrentTask", m.CurrentTask)),
				NotItem(EqualItem("Tasks", m.Tasks)),
				NotItem(EqualItem("Stage", m.Stage)))))
}

func (a *MachineAgent) ChangeStage() {
	var cmObj interface{}
	a.state = AGENT_WAIT_FOR_RUNNABLE
	if strings.HasSuffix(a.machine.BootEnv, "-install") {
		a.state = AGENT_EXIT
	}
	csMap := map[string]string{}
	nextStage := ""
	csErr := a.client.Req().Get().
		UrlForM(a.machine, "params", "change-stage/map").
		Params("aggregate", "true").Do(&cmObj)
	if csErr != nil {
		// No change stage map, see if we need to go to local
		if strings.HasSuffix(a.machine.BootEnv, "-install") {
			nextStage = "local"
		}
	} else {
		if err := utils.Remarshal(cmObj, &csMap); err != nil {
			a.err = err
			a.initOrExit()
			return
		}
		ns, ok := csMap[a.machine.Stage]
		if !ok {
			return
		}
		pieces := strings.SplitN(ns, ":", 2)
		nextStage = pieces[0]
		if len(pieces) == 2 {
			switch pieces[1] {
			case "Reboot":
				a.state = AGENT_REBOOT
			case "Stop":
				a.state = AGENT_EXIT
			case "Shutdown":
				a.state = AGENT_POWEROFF
			}
		}
	}
	if nextStage != "" {
		newM := models.Clone(a.machine).(*models.Machine)
		newM.Stage = nextStage
		if _, err := a.client.PatchTo(a.machine, newM); err != nil {
			a.err = err
			a.initOrExit()
		}
	}
}

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
			a.Init()
		case AGENT_WAIT_FOR_RUNNABLE:
			a.WaitRunnable()
		case AGENT_RUN_TASK:
			a.RunTask()
		case AGENT_WAIT_FOR_CHANGE_STAGE:
			a.WaitChangeStage()
		case AGENT_CHANGE_STAGE:
			a.ChangeStage()
		case AGENT_EXIT:
			return a.err
		case AGENT_REBOOT:
			return a.power("reboot")
		case AGENT_POWEROFF:
			return a.power("poweroff")
		}
	}
}
