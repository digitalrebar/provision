package models

import (
	"fmt"
	"time"

	"github.com/pborman/uuid"
)

// AsyncAction is custom workflow-like element that defines parameters and an AsyncActionTemplate
// that acts as a basis for expanding into tasks.
//
// swagger:model
type AsyncAction struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	// Uuid is the key of this particular AsyncAction.
	// required: true
	// swagger:strfmt uuid
	Uuid uuid.UUID `index:",key"`
	// AsyncActionTemplate defines the tasks and base parameters for this action
	// required: true
	AsyncActionTemplate string
	// State The state the async action is in.  Must be one of "created", "running", "failed", "finished", "incomplete", "cancelled"
	// required: true
	State string
	// StartTime The time the async action started running.
	StartTime time.Time
	// EndTime The time the async action failed or finished or cancelled.
	EndTime time.Time
	// Machine is the key of the machine running the AsyncAction
	// swagger:strfmt uuid
	Machine uuid.UUID
	// Current indicates that this is the current action
	Current bool
	// Previous in the chain for this machine
	// swagger:strfmt uuid
	Previous uuid.UUID
	// Next in the chain for this machine
	// swagger:strfmt uuid
	Next uuid.UUID
	// Archived indicates whether the complete log for the async action can be
	// retrieved via the API.  If Archived is true, then the log cannot
	// be retrieved.
	//
	// required: true
	Archived bool
	// An array of profiles to apply to this machine in order when looking
	// for a parameter during rendering.
	Profiles []string
	// The Parameters that have been directly set on the Machine.
	Params map[string]interface{}
}

func (aa *AsyncAction) GetMeta() Meta {
	return aa.Meta
}

func (aa *AsyncAction) SetMeta(d Meta) {
	aa.Meta = d
}

func (aa *AsyncAction) Validate() {
	aa.AddError(ValidName("Invalid AsyncActionTemplate", aa.AsyncActionTemplate))
	switch aa.State {
	case "created", "running", "incomplete":
	case "failed", "finished", "cancelled":
	default:
		aa.AddError(fmt.Errorf("Invalid State `%s`", aa.State))
	}
}

func (aa *AsyncAction) Prefix() string {
	return "async_actions"
}

func (aa *AsyncAction) Key() string {
	return aa.Uuid.String()
}

func (aa *AsyncAction) KeyName() string {
	return "Uuid"
}

func (aa *AsyncAction) Fill() {
	if aa.Meta == nil {
		aa.Meta = Meta{}
	}
	if aa.Profiles == nil {
		aa.Profiles = []string{}
	}
	if aa.Params == nil {
		aa.Params = map[string]interface{}{}
	}
	aa.Validation.fill(aa)
}

func (aa *AsyncAction) AuthKey() string {
	return aa.Machine.String()
}

func (aa *AsyncAction) SliceOf() interface{} {
	s := []*AsyncAction{}
	return &s
}

func (aa *AsyncAction) ToModels(obj interface{}) []Model {
	items := obj.(*[]*AsyncAction)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (aa *AsyncAction) CanHaveActions() bool {
	return true
}

// match Profiler interface

// GetProfiles gets the profiles on this stage
func (aa *AsyncAction) GetProfiles() []string {
	return aa.Profiles
}

// SetProfiles sets the profiles on this stage
func (aa *AsyncAction) SetProfiles(p []string) {
	aa.Profiles = p
}

// match Paramer interface

// GetParams gets the parameters on this stage
func (aa *AsyncAction) GetParams() map[string]interface{} {
	return copyMap(aa.Params)
}

// SetParams sets the parameters on this stage
func (aa *AsyncAction) SetParams(p map[string]interface{}) {
	aa.Params = copyMap(p)
}
