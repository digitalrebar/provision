package api

import (
	"testing"
	"time"

	"github.com/digitalrebar/provision/models"
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
	recieved := <-ch
	t.Logf("Recieved event: %#v", recieved)
}
