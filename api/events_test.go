package api

import (
	"syscall"
	"testing"
	"time"

	"github.com/digitalrebar/provision/v4/models"
)

func TestEvents(t *testing.T) {
	t.Logf("Setting up event listener")
	listener, err := session.Events()
	if err != nil {
		t.Errorf("Failed to create EventStream: %v", err)
		return
	}
	t.Logf("Listening for beeblebrox events")
	handle, ch, err := listener.Register("beeblebrox.*.*")
	defer listener.Deregister(handle)
	if err != nil {
		t.Errorf("Failed to register for beebleborx events: %v", err)
		return
	}
	evt := &models.Event{
		Time:   time.Now(),
		Type:   "beeblebrox",
		Action: "created",
		Key:    "foo",
	}
	if err := session.PostEvent(evt); err != nil {
		t.Errorf("Failed to create new Event: %v", err)
		return
	}
	t.Logf("Waiting for event from server")
	received := <-ch
	t.Logf("Received event: %#v", received)
}

func TestEventDeadlock(t *testing.T) {
	t.Logf("Setting up event listener to deadlock")
	listener, err := session.Events()
	if err != nil {
		t.Errorf("Failed to create EventStream: %v", err)
		return
	}
	t.Logf("Listening for users events")
	handle, ch, err := listener.Register("users.*.*")
	defer listener.Deregister(handle)
	if err != nil {
		t.Errorf("Failed to register for users events: %v", err)
		return
	}

	done := make(chan bool)
	finished := make(chan bool)
	go func() {
		leave := false
		for !leave {
			select {
			case <-done:
				leave = true
			default:
				user := &models.User{Name: "user1"}
				session.FillModel(user, "user1")
			}
		}
		finished <- true
	}()

	user1 := &models.User{Name: "user1"}
	if err := session.CreateModel(user1); err != nil {
		t.Errorf("Failed to create user1 for users events: %v", err)
		return
	}
	if _, err := session.DeleteModel("users", user1.Name); err != nil {
		t.Errorf("Failed to destroy user1 for users events: %v", err)
		return
	}

	t.Logf("Waiting for event from server")
	received := <-ch
	t.Logf("Received event: %#v", received)
	received = <-ch
	t.Logf("Received event: %#v", received)
	done <- true
	<-finished
}

func TestWaitFor(t *testing.T) {
	machine1 := mustDecode(&models.Machine{}, `
Address: 192.168.100.110
Endpoint: Fred
BootEnv: local
Meta:
  feature-flags: change-stage-v2
Name: john
Uuid: 3e7031fe-3062-45f1-835c-92541bc9cbd3
Validated: true
`).(*models.Machine)

	machineRes := models.Clone(machine1).(*models.Machine)
	machineRes.Secret = ""
	machineRes.Runnable = true
	machineRes.Stage = "none"
	machineRes.CurrentTask = 0
	machineRes.Pool = "default"
	machineRes.PoolStatus = "Free"
	machineRes.WorkflowComplete = true
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

	es, err := session.Events()
	if err != nil {
		t.Errorf("Failed to create new EventStream: %v", err)
		return
	}
	defer es.Close()

	answer, err := es.WaitFor(machine1, EqualItem("Runnable", false), 1000)
	if err != nil {
		t.Errorf("WaitFor should not have returned an error: %v", err)
		return
	}
	if answer != "timeout" {
		t.Errorf("WaitFor should have timed out: %s", answer)
	}

	answer, err = es.WaitFor(machine1, EqualItem("Runnable", true), 1000)
	if err != nil {
		t.Errorf("WaitFor should not have returned an error: %v", err)
		return
	}
	if answer != "complete" {
		t.Errorf("WaitFor should have completed: %s", answer)
	}

	mm1 := models.Clone(machine1).(*models.Machine)
	go func() {
		mc := models.Clone(mm1).(*models.Machine)
		time.Sleep(1 * time.Second)
		mc.Runnable = false
		_, err := session.PatchTo(mm1, mc)
		if err != nil {
			t.Errorf("Failed to update runnable: %v", err)
		}
	}()

	answer, err = es.WaitFor(machine1, EqualItem("Runnable", false), 3*time.Second)
	if err != nil {
		t.Errorf("WaitFor should not have returned an error: %v", err)
		return
	}
	if answer != "complete" {
		t.Errorf("WaitFor should have completed: %s", answer)
	}

	go func() {
		time.Sleep(1 * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()

	answer, err = es.WaitFor(machine1, EqualItem("Runnable", true), 3*time.Second)
	if err != nil {
		t.Errorf("WaitFor should not have returned an error: %v", err)
		return
	}
	if answer != "interrupt" {
		t.Errorf("WaitFor should have interrupt: %s", answer)
	}

	session.Req().Delete(machine1)
}
