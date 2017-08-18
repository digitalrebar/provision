package models

import (
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

	// Object - the data of the object.
	Object interface{}
}

func (e *Event) Text() string {
	jsonString, err := json.MarshalIndent(e.Object, "", "  ")
	if err != nil {
		jsonString = []byte("json failure")
	}

	return fmt.Sprintf("%d: %s %s %s\n%s\n", e.Time.Unix(), e.Type, e.Action, e.Key, string(jsonString))
}
