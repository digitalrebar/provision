package api

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/digitalrebar/provision/models"
)

func runAgent(t *testing.T, mi *models.Machine, lastTask, lastState, lastExitState string) (m *models.Machine) {
	t.Helper()
	m = mi

	if !m.Runnable {
		t.Logf("Machine %s not runnable, patching it to a Runnable state", m.Name)
		mc := models.Clone(m).(*models.Machine)
		mc.Runnable = true
		mr, err := session.PatchTo(m, mc)
		if err != nil {
			t.Errorf("ERROR: Failed to make machine runnable: %v", err)
			return
		}
		m = mr.(*models.Machine)
		t.Logf("Machine %s patched", m.Name)
	}

	buf := &bytes.Buffer{}
	err := session.Agent(m, true, true, false, buf)
	if err != nil {
		t.Errorf("ERROR: Agent run failed: %v", err)
		return
	}
	if err := session.FillModel(m, m.Key()); err != nil {
		t.Errorf("ERROR: Failed to fetch machine1: %v", err)
		return
	}
	t.Logf("Machine current job: %v", m.CurrentJob)
	job := &models.Job{Uuid: m.CurrentJob}
	if err := session.FillModel(job, job.Key()); err != nil {
		t.Errorf("ERROR: Failed to fetch current job: %v", err)
		return
	}
	t.Logf("Job log: \n---------------\n")
	t.Logf("%s", buf.String())
	t.Logf("\n---------------\nEnd log\n\n")
	if job.Task == lastTask && job.State == lastState && job.ExitState == lastExitState {
		t.Logf("Run for task %s finished with desired state %s:%s", job.Task, job.State, job.ExitState)
	} else {
		t.Errorf("ERROR: Run for task %s finished with unknown state %s:%s", job.Task, job.State, job.ExitState)
	}
	return
}

func TestJobs(t *testing.T) {
	tjd, err := ioutil.TempDir("", "jobTest-")
	if err != nil {
		t.Errorf("Failed to create tmpdir for job tester")
		return
	}
	defer os.RemoveAll(tjd)
	os.Setenv("JT", tjd)
	machine1 := mustDecode(&models.Machine{}, `
Address: 192.168.100.110
BootEnv: local
Meta:
  feature-flags: change-stage-v2
Name: john
Uuid: 3e7031fe-3062-45f1-835c-92541bc9cbd3
Validated: true
`).(*models.Machine)
	task1 := mustDecode(&models.Task{}, `
Name: task1
Meta:
  feature-flags: original-exit-codes
Templates:
  - Name: reboot
    Contents: |
      #!/usr/bin/env bash
      [[ -e "$JT"/reboot-orig.txt ]] && exit 0
      touch "$JT"/reboot-orig.txt""
      exit 1
  - Name: fail
    Contents: |
      #!/usr/bin/env bash
      [[ -e "$JT"/fail-orig.txt ]] && exit 0
      touch "$JT"/fail-orig.txt
      exit 4
  - Name: incomplete
    Contents: |
      #!/usr/bin/env bash
      [[ -e "$JT"/incomplete-orig.txt ]] && exit 0
      touch "$JT"/incomplete-orig.txt
      exit 2
`).(*models.Task)

	task2 := mustDecode(&models.Task{}, `
Name: task2
Meta:
  feature-flags: sane-exit-codes
Templates:
  - Name: reboot
    Contents: |
      #!/usr/bin/env bash
      [[ -e "$JT"/reboot.txt ]] && exit 0
      touch "$JT"/reboot.txt
      exit 64
  - Name: poweroff
    Contents: |
      #!/usr/bin/env bash
      [[ -e "$JT"/poweroff.txt ]] && exit 0
      touch "$JT"/poweroff.txt
      exit 32
  - Name: stop
    Contents: |
      #!/usr/bin/env bash
      [[ -e "$JT"/stop.txt ]] && exit 0
      touch "$JT"/stop.txt
      exit 16
  - Name: fail
    Contents: |
      #!/usr/bin/env bash
      [[ -e "$JT"/fail.txt ]] && exit 0
      touch "$JT"/fail.txt
      exit 1
  - Name: incomplete
    Contents: |
      #!/usr/bin/env bash
      [[ -e "$JT"/incomplete.txt ]] && exit 0
      touch "$JT"/incomplete.txt
      exit 128
`).(*models.Task)

	stage1 := mustDecode(&models.Stage{}, `
Name: stage1
Tasks:
- task1
- task2
`).(*models.Stage)
	stage2 := mustDecode(&models.Stage{}, `
Name: stage2
Tasks:
- task2
- task1
`).(*models.Stage)
	machineRes := models.Clone(machine1).(*models.Machine)
	machineRes.Secret = ""
	machineRes.Runnable = true
	machineRes.Stage = "none"
	rt(t, "Make initial machine", machineRes, nil,
		func() (interface{}, error) {
			err := session.CreateModel(machine1)
			if err != nil {
				return machine1, err
			}
			res := models.Clone(machine1).(*models.Machine)
			res.Secret = ""
			return res, err
		}, nil)
	rt(t, "Make task1", models.Clone(task1), nil,
		func() (interface{}, error) {
			err := session.CreateModel(task1)
			return task1, err
		}, nil)
	rt(t, "Make task2", models.Clone(task2), nil,
		func() (interface{}, error) {
			err := session.CreateModel(task2)
			return task2, err
		}, nil)
	rt(t, "Make stage1", models.Clone(stage1), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stage1)
			return stage1, err
		}, nil)
	rt(t, "Make stage2", models.Clone(stage2), nil,
		func() (interface{}, error) {
			err := session.CreateModel(stage2)
			return stage2, err
		}, nil)
	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stage1"
	machineRes.Tasks = []string{"task1", "task2"}
	machineRes.CurrentTask = -1
	rt(t, "Set machine 1 to stage1", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stage1"
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	rt(t, "Increment machine1's CurrentTask pointer", nil,
		&models.Error{
			Model:    "machines",
			Key:      "3e7031fe-3062-45f1-835c-92541bc9cbd3",
			Type:     "ValidationError",
			Messages: []string{"Cannot change CurrentTask from -1 to 0"},
			Code:     422,
		},
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.CurrentTask += 1
			return session.PatchTo(machine1, mc)
		}, nil)
	rt(t, "Set machine 1 to stage2 (without force)", nil,
		&models.Error{
			Model:    "machines",
			Key:      machine1.UUID(),
			Type:     "ValidationError",
			Messages: []string{"Can not change stages with pending tasks unless forced"},
			Code:     422,
		},
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stage2"
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Tasks = []string{}
	rt(t, "Try to remove tasks from machine1", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Tasks = []string{"task2", "task1"}
	rt(t, "Try to change order of tasks on a machine", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{"task2", "task1"}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Tasks = []string{"task1", "task2", "task1"}
	rt(t, "Append a task to a machine", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{"task1", "task2", "task1"}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Stage = "stage2"
	machineRes.Tasks = []string{"task2", "task1"}
	machineRes.CurrentTask = -1
	rt(t, "Set machine 1 to stage2 (forced)", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Stage = "stage2"
			err := session.Req().PatchTo(machine1, mc).Params("force", "true").Do(&mc)
			if err == nil {
				machine1 = mc
			}
			return mc, err
		}, nil)
	machine1 = runAgent(t, machine1, "task2", "incomplete", "reboot")
	machine1 = runAgent(t, machine1, "task2", "incomplete", "poweroff")
	machine1 = runAgent(t, machine1, "task2", "incomplete", "stop")
	machine1 = runAgent(t, machine1, "task2", "failed", "failed")
	rt(t, "Try to remove tasks from machine1 (fail)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "3e7031fe-3062-45f1-835c-92541bc9cbd3",
			Type:     "ValidationError",
			Messages: []string{"Cannot remove tasks that have already executed or are already executing"},
			Code:     422,
		},
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	rt(t, "Try to change order of tasks on a machine (fail)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "3e7031fe-3062-45f1-835c-92541bc9cbd3",
			Type:     "ValidationError",
			Messages: []string{"Cannot change tasks that have already executed or are executing"},
			Code:     422,
		},
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{"task1", "task2"}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Tasks = []string{"task2"}
	rt(t, "Remove a to-be-executed task", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{"task2"}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Tasks = []string{"task2", "task1", "task2"}
	rt(t, "Append extra tasks in the middle of a run", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{"task2", "task1", "task2"}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	machine1 = runAgent(t, machine1, "task2", "incomplete", "complete")
	machine1 = runAgent(t, machine1, "task1", "incomplete", "reboot")
	machine1 = runAgent(t, machine1, "task1", "failed", "failed")
	machine1 = runAgent(t, machine1, "task1", "incomplete", "complete")
	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.Tasks = []string{"task2", "task1"}
	rt(t, "Remove extra task in the middle of a run", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{"task2", "task1"}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	rt(t, "Try to remove tasks from machine1 (fail again)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "3e7031fe-3062-45f1-835c-92541bc9cbd3",
			Type:     "ValidationError",
			Messages: []string{"Cannot remove tasks that have already executed or are already executing"},
			Code:     422,
		},
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{}
			return session.PatchTo(machine1, mc)
		}, nil)
	rt(t, "Try to change order of tasks on a machine (fail again)", nil,
		&models.Error{
			Model:    "machines",
			Key:      "3e7031fe-3062-45f1-835c-92541bc9cbd3",
			Type:     "ValidationError",
			Messages: []string{"Cannot change tasks that have already executed or are executing"},
			Code:     422,
		},
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{"task1", "task2"}
			return session.PatchTo(machine1, mc)
		}, nil)
	machine1 = runAgent(t, machine1, "task1", "finished", "complete")
	if machine1.CurrentTask == len(machine1.Tasks) {
		t.Logf("All tasks on machine1 finished")
	} else {
		t.Errorf("ERROR: Machine1: currentTask %d, tasks %v:%d", machine1.CurrentTask, machine1.Tasks, len(machine1.Tasks))
	}
	machineRes = models.Clone(machine1).(*models.Machine)
	machineRes.CurrentTask = -1
	rt(t, "Increment machine1's CurrentTask pointer", machineRes, nil,
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.CurrentTask = -1
			return session.PatchTo(machine1, mc)
		}, nil)
	j := []*models.Job{}
	if err := session.Req().Filter("jobs",
		"Current", "Eq", "true").
		Do(&j); err != nil {
		t.Errorf("Error getting jobs: %v", err)
	} else if len(j) != 1 {
		t.Errorf("Expected 1 current job, not %d", len(j))
	} else if j[0].Key() != machine1.CurrentJob.String() {
		t.Errorf("Expected current job to match what was recorded on machine1, not %s", j[0].Key())
	} else {
		t.Logf("Got expected current job results")
	}
	session.Req().Delete(machine1)
	session.Req().Delete(stage1)
	session.Req().Delete(stage2)
	session.Req().Delete(task1)
	session.Req().Delete(task2)
	j = []*models.Job{}
	if err := session.Req().UrlFor("jobs").Do(&j); err != nil {
		t.Errorf("Error getting jobs: %v", err)
	} else if len(j) != 4 {
		t.Errorf("Expected 4 jobs, not %d", len(j))
	} else {
		t.Logf("Got expected 4 jobs")
		for _, job := range j {
			session.Req().Delete(job)
		}
	}

}
