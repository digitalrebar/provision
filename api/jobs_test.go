package api

import (
	"testing"

	"github.com/digitalrebar/provision/models"
)

func TestJobs(t *testing.T) {
	machine1 := mustDecode(&models.Machine{}, `
Address: 192.168.100.110
BootEnv: local
Name: john
Uuid: 3e7031fe-3062-45f1-835c-92541bc9cbd3
Validated: true
`).(*models.Machine)
	task1 := mustDecode(&models.Task{}, `
Name: task1
Meta:
  feature-flags: original-exit-codes
Templates:
  - Name: expando
    Path: /tmp/expando.txt
    Contents: "Hey I am a test content"
  - Name: reboot
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/reboot.txt ]] && exit 0
      touch /tmp/reboot.txt""
      exit 1
  - Name: incomplete
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/incomplete.txt ]] && exit 0
      touch /tmp/incomplete.txt
      exit 2
  - Name: incomplete-reboot
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/i-r.txt ]] && exit 0
      touch /tmp/i-r.txt
      exit 3
  - Name: fail
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/fail.txt ]] && exit 0
      touch /tmp/fail.txt
      exit 4
  - Name: success
    Contents: |
      #!/usr/bin/env true
`).(*models.Task)

	task2 := mustDecode(&models.Task{}, `
Name: task2
Meta:
  feature-flags: sane-exit-codes
Templates:
  - Name: expando
    Path: /tmp/expando.txt
    Contents: "Hey I am a test content"
  - Name: reboot
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/reboot.txt ]] && exit 0
      touch /tmp/reboot.txt""
      exit 64
  - Name: incomplete
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/incomplete.txt ]] && exit 0
      touch /tmp/incomplete.txt
      exit 128
  - Name: incomplete-reboot
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/i-r.txt ]] && exit 0
      touch /tmp/i-r.txt
      exit 192
  - Name: poweroff
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/poweroff.txt ]] && exit 0
      touch /tmp/poweroff.txt
      exit 32
  - Name: stop
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/stop.txt ]] && exit 0
      touch /tmp/stop.txt
      exit 16
  - Name: fail
    Contents: |
      #!/usr/bin/env bash
      [[ -e /tmp/fail.txt ]] && exit 0
      touch /tmp/fail.txt
      exit 1
  - Name: success
    Contents: |
      #!/usr/bin/env true
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
	rt(t, "Try to remove tasks from machine1", nil,
		&models.Error{
			Model: "machines",
			Key:   machine1.UUID(),
			Type:  "ValidationError",
			Messages: []string{
				"Cannot remove tasks from machines without changing stage",
				"Can only append tasks to the task list on a machine."},
			Code: 422,
		},
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{"task2"}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	rt(t, "Try to change order of tasks on a machine", nil,
		&models.Error{
			Model:    "machines",
			Key:      machine1.UUID(),
			Type:     "ValidationError",
			Messages: []string{"Can only append tasks to the task list on a machine."},
			Code:     422,
		},
		func() (interface{}, error) {
			mc := models.Clone(machine1).(*models.Machine)
			mc.Tasks = []string{"task2", "task1"}
			res, err := session.PatchTo(machine1, mc)
			if err == nil {
				machine1 = res.(*models.Machine)
			}
			return res, err
		}, nil)
	rt(t, "Try to change order of tasks on a machine", nil,
		&models.Error{
			Model:    "machines",
			Key:      machine1.UUID(),
			Type:     "ValidationError",
			Messages: []string{"Can only append tasks to the task list on a machine."},
			Code:     422,
		},
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
}
