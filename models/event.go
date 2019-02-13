package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Event represents an action in the system.
// In general, the event generates for a subject
// of the form: type.action.key
//
// swagger:model
type Event struct {
	// Time of the event.
	// swagger:strfmt date-time
	Time time.Time

	// Type - object type
	Type string

	// Action - what happened
	Action string

	// Key - the id of the object
	Key string

	// Principal - the user or subsystem that caused the event to be emitted
	Principal string

	// Object - the data of the object.
	Object interface{}

	// Original - the data of the object before the operation (update and save only)
	Original interface{}
}

func (e *Event) Text() string {
	jsonString, err := json.MarshalIndent(e.Object, "", "  ")
	if err != nil {
		jsonString = []byte("json failure")
	}

	return fmt.Sprintf("%d: %s %s %s %s\n%s\n", e.Time.Unix(), e.Type, e.Action, e.Key, e.Principal, string(jsonString))
}

func (e *Event) Model() (Model, error) {
	res, err := New(e.Type)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	enc, dec := json.NewEncoder(&buf), json.NewDecoder(&buf)
	err = enc.Encode(e.Object)
	if err != nil {
		return nil, err
	}
	err = dec.Decode(res)
	return res, err
}

func (e *Event) Message() string {
	if s, ok := e.Object.(string); ok {
		return s
	}
	return ""
}

func EventFor(obj Model, action string) *Event {
	return &Event{
		Time:   time.Now(),
		Type:   obj.Prefix(),
		Action: action,
		Key:    obj.Key(),
		Object: obj,
	}
}
