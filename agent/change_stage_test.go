package agent

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
)

func TestChangeStage(t *testing.T) {
	tjd, err := ioutil.TempDir("", "changeStage-")
	if err != nil {
		t.Errorf("Failed to create tmpdir for change stage tester")
		return
	}
	defer os.RemoveAll(tjd)
	os.Setenv("JT", tjd)
	if session == nil {
		session, err = api.UserSession("https://127.0.0.1:10021", "rocketskates", "r0cketsk8ts")
		if err != nil {
			t.Errorf("Error creating session: %v", err)
			return
		}
		defer func() { session = nil }()
	}

	machine1 := mustDecode(&models.Machine{}, `
Address: 192.168.100.110
Endpoint: ""
BootEnv: local
Meta:
  feature-flags: ""
Name: john
Params:
  change-stage/map:
    stageNoWait: stageDoneNoWait:Success
    stageNoWait1: stageDoneWait:Success
    stageWait: stageDoneNoWait:Success
    stageWait1: stageDoneWait:Success
    stageStop: stageDoneWait:Stop
    stageReboot: stageDoneWait:Reboot
    stageRealReboot: stageDoneWait:Success
    stageRealReboot1: stageDoneReboot:Success
    fred-install: stageDoneWait:Success
Uuid: 3e7031fe-3062-45f1-835c-92541bc9cbd3
Validated: true
`).(*models.Machine)

	taskInstall := mustDecode(&models.Task{}, `
Name: taskInstall
Endpoint: ""
Meta:
  feature-flags: sane-exit-codes
Templates:
  - Name: success
    Contents: |
      #!/usr/bin/env bash
      exit 0
    Meta: {}
`).(*models.Task)

	task1 := mustDecode(&models.Task{}, `
Name: task1
Endpoint: ""
Meta:
  feature-flags: sane-exit-codes
Templates:
  - Name: success
    Contents: |
      #!/usr/bin/env bash
      exit 0
    Meta: {}
`).(*models.Task)

	taskDone := mustDecode(&models.Task{}, `
Name: taskDone
Endpoint: ""
Meta:
  feature-flags: sane-exit-codes
Templates:
  - Name: success
    Contents: |
      #!/usr/bin/env bash
      exit 0
    Meta: {}
`).(*models.Task)

	stageGregInstall := mustDecode(&models.Stage{}, `
Name: greg-install
Endpoint: ""
RunnerWait: true
Tasks:
- taskInstall
`).(*models.Stage)
	stageFredInstall := mustDecode(&models.Stage{}, `
Name: fred-install
Endpoint: ""
RunnerWait: true
Tasks:
- taskInstall
`).(*models.Stage)
	stageRealReboot1 := mustDecode(&models.Stage{}, `
Name: stageRealReboot1
Endpoint: ""
RunnerWait: true
Tasks:
- task1
`).(*models.Stage)
	stageRealReboot := mustDecode(&models.Stage{}, `
Name: stageRealReboot
Endpoint: ""
RunnerWait: true
Tasks:
- task1
`).(*models.Stage)
	stageStop := mustDecode(&models.Stage{}, `
Name: stageStop
Endpoint: ""
RunnerWait: true
Tasks:
- task1
`).(*models.Stage)
	stageReboot := mustDecode(&models.Stage{}, `
Name: stageReboot
Endpoint: ""
RunnerWait: true
Tasks:
- task1
`).(*models.Stage)
	stageNoWait := mustDecode(&models.Stage{}, `
Name: stageNoWait
Endpoint: ""
RunnerWait: true
Tasks:
- task1
`).(*models.Stage)
	stageWait := mustDecode(&models.Stage{}, `
Name: stageWait
Endpoint: ""
RunnerWait: true
Tasks:
- task1
`).(*models.Stage)
	stageNoWait1 := mustDecode(&models.Stage{}, `
Name: stageNoWait1
Endpoint: ""
RunnerWait: true
Tasks:
- task1
`).(*models.Stage)
	stageWait1 := mustDecode(&models.Stage{}, `
Name: stageWait1
Endpoint: ""
RunnerWait: true
Tasks:
- task1
`).(*models.Stage)
	stageDoneNoWait := mustDecode(&models.Stage{}, `
Name: stageDoneNoWait
Endpoint: ""
RunnerWait: true
Tasks:
- taskDone
`).(*models.Stage)
	stageDoneWait := mustDecode(&models.Stage{}, `
Name: stageDoneWait
Endpoint: ""
RunnerWait: true
Tasks:
- taskDone
`).(*models.Stage)
	stageDoneReboot := mustDecode(&models.Stage{}, `
Name: stageDoneReboot
Endpoint: ""
RunnerWait: true
Reboot: true
Tasks:
- taskDone
`).(*models.Stage)

	machineRes := models.Clone(machine1).(*models.Machine)
	machineRes.Secret = ""
	machineRes.Runnable = true
	machineRes.Stage = "none"
	machineRes.CurrentTask = 0
	machineRes.WorkflowComplete = true
	machineRes.Pool = "default"
	machineRes.PoolStatus = "Free"
	rt(t, "Make initial machine", machineRes, nil,
		func() (interface{}, error) {
			err := session.CreateModel(machine1)
			if err != nil {
				return machine1, err
			}
			res := models.Clone(machine1).(*models.Machine)
			res.Secret = ""
			res.Meta["feature-flags"] = ""
			return res, err
		}, nil)
	rt(t, "Make task1", models.Clone(task1), nil,
		func() (interface{}, error) {
			err := session.CreateModel(task1)
			return task1, err
		}, nil)
	rt(t, "Make taskInstall", models.Clone(taskInstall), nil,
		func() (interface{}, error) {
			err := session.CreateModel(taskInstall)
			return taskInstall, err
		}, nil)
	rt(t, "Make taskDone", models.Clone(taskDone), nil,
		func() (interface{}, error) {
			err := session.CreateModel(taskDone)
			return taskDone, err
		}, nil)
	rt(t, "Make stageNoWait", models.Clone(stageNoWait), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageNoWait)
			return stageNoWait, err
		}, nil)
	rt(t, "Make stageWait", models.Clone(stageWait), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageWait)
			return stageWait, err
		}, nil)
	rt(t, "Make stageNoWait1", models.Clone(stageNoWait1), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageNoWait1)
			return stageNoWait1, err
		}, nil)
	rt(t, "Make stageWait1", models.Clone(stageWait1), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageWait1)
			return stageWait1, err
		}, nil)
	rt(t, "Make stageDoneNoWait", models.Clone(stageDoneNoWait), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageDoneNoWait)
			return stageDoneNoWait, err
		}, nil)
	rt(t, "Make stageDoneWait", models.Clone(stageDoneWait), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageDoneWait)
			return stageDoneWait, err
		}, nil)
	rt(t, "Make stageDoneReboot", models.Clone(stageDoneReboot), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageDoneReboot)
			return stageDoneReboot, err
		}, nil)
	rt(t, "Make stageReboot", models.Clone(stageReboot), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageReboot)
			return stageReboot, err
		}, nil)
	rt(t, "Make stageGregInstall", models.Clone(stageGregInstall), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageGregInstall)
			return stageGregInstall, err
		}, nil)
	rt(t, "Make stageFredInstall", models.Clone(stageFredInstall), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageFredInstall)
			return stageFredInstall, err
		}, nil)
	rt(t, "Make stageRealReboot", models.Clone(stageRealReboot), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageRealReboot)
			return stageRealReboot, err
		}, nil)
	rt(t, "Make stageRealReboot1", models.Clone(stageRealReboot1), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageRealReboot1)
			return stageRealReboot1, err
		}, nil)
	rt(t, "Make stageStop", models.Clone(stageStop), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stageStop)
			return stageStop, err
		}, nil)

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stageNoWait"
	machineRes.Tasks = []string{"task1"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to stageNoWait", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stageNoWait"
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "taskDone", "finished", "complete")

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stageWait"
	machineRes.Tasks = []string{"task1"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to stageWait", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stageWait"
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "taskDone", "finished", "complete")

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stageNoWait1"
	machineRes.Tasks = []string{"task1"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to stageNoWait1", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stageNoWait1"
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "taskDone", "finished", "complete")

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stageWait1"
	machineRes.Tasks = []string{"task1"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to stageWait1", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stageWait1"
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "taskDone", "finished", "complete")

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stageReboot"
	machineRes.Tasks = []string{"task1"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to stageReboot", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stageReboot"
			res := models.Clone(machine1)
			err := session.Req().PatchTo(machine1, mc).Params("force", "true").Do(&res)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "task1", "finished", "complete")

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stageRealReboot"
	machineRes.Tasks = []string{"task1"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to stageRealReboot", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stageRealReboot"
			res := models.Clone(machine1)
			err := session.Req().PatchTo(machine1, mc).Params("force", "true").Do(&res)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "taskDone", "finished", "complete")

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stageRealReboot1"
	machineRes.Tasks = []string{"task1"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to stageRealReboot1", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stageRealReboot1"
			res := models.Clone(machine1)
			err := session.Req().PatchTo(machine1, mc).Params("force", "true").Do(&res)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "task1", "finished", "complete")

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stageStop"
	machineRes.Tasks = []string{"task1"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to stageStop", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stageStop"
			res := models.Clone(machine1)
			err := session.Req().PatchTo(machine1, mc).Params("force", "true").Do(&res)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "task1", "finished", "complete")

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "fred-install"
	machineRes.Tasks = []string{"taskInstall"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to fred-install", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "fred-install"
			res := models.Clone(machine1)
			err := session.Req().PatchTo(machine1, mc).Params("force", "true").Do(&res)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "taskDone", "finished", "complete")

	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "greg-install"
	machineRes.Tasks = []string{"taskInstall"}
	machineRes.CurrentTask = -1
	machineRes.WorkflowComplete = false
	rt(t, "Set machine 1 to greg-install", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "greg-install"
			res := models.Clone(machine1)
			err := session.Req().PatchTo(machine1, mc).Params("force", "true").Do(&res)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "taskInstall", "finished", "complete")

	session.Req().Delete(machine1)
	session.Req().Delete(stageFredInstall)
	session.Req().Delete(stageGregInstall)
	session.Req().Delete(stageReboot)
	session.Req().Delete(stageRealReboot)
	session.Req().Delete(stageRealReboot1)
	session.Req().Delete(stageStop)
	session.Req().Delete(stageNoWait)
	session.Req().Delete(stageWait)
	session.Req().Delete(stageNoWait1)
	session.Req().Delete(stageWait1)
	session.Req().Delete(stageDoneNoWait)
	session.Req().Delete(stageDoneWait)
	session.Req().Delete(stageDoneReboot)
	session.Req().Delete(task1)
	session.Req().Delete(taskInstall)
	session.Req().Delete(taskDone)
	j := []*models.Job{}
	if err := session.Req().UrlFor("jobs").Do(&j); err != nil {
		t.Errorf("Error getting jobs: %v", err)
	} else {
		jobCount := 16
		if len(j) != jobCount {
			t.Errorf("Expected %d jobs, not %d", jobCount, len(j))
		} else {
			t.Logf("Got expected %d jobs", jobCount)
		}
		for _, job := range j {
			session.Req().Delete(job)
		}
	}
}
