package models

import (
	"fmt"
)

/*
 * Endpoint tracks who to get events from and where to send them when
 * actions are required.
 */

// Element define a part of the endpoint
type Element struct {
	Type          string
	Version       string
	Name          string
	ActualVersion string
}

// ElementAction defines an action to take on an Element
type ElementAction struct {
	Element
	Action string
	Value  interface{}
}

func (ea *ElementAction) String() string {
	return fmt.Sprintf("%s %s %s:%s", ea.Action, ea.Type, ea.Name, ea.Version)
}

// Endpoint represents a managed Endpoint
// swagger:model
type Endpoint struct {
	Validation
	Access
	Meta
	Owned
	Bundled

	Id string

	Description   string `json:"Description,omitempty"`
	Documentation string `json:"Documentation,omitempty"`

	// Holds the access parameters.
	Params map[string]interface{} `json:"Params,omitempty"`

	ConnectionStatus string `json:"ConnectionStatus,omitempty"`

	// Deprecated
	VersionSet string `json:"VersionSet,omitempty"`

	// VersionSets replaces VersionSet - code processes both
	VersionSets []string `json:"VersionSets,omitempty"`

	Apply        bool                   `json:"Apply,omitempty"`
	HaId         string                 `json:"HaId,omitempty"`
	Arch         string                 `json:"Arch,omitempty"`
	Os           string                 `json:"Os,omitempty"`
	DRPVersion   string                 `json:"DRPVersion,omitempty"`
	DRPUXVersion string                 `json:"DRPUXVersion,omitempty"`
	Components   []*Element             `json:"Components,omitempty"`
	Plugins      []*Plugin              `json:"Plugins,omitempty"`
	Prefs        map[string]string      `json:"Prefs,omitempty"`
	Global       map[string]interface{} `json:"Global,omitempty"`
	Actions      []*ElementAction       `json:"Actions,omitempty"`
}

// GetMeta get the meta data from the model
func (e *Endpoint) GetMeta() Meta {
	return e.Meta
}

// SetMeta set the meta data on the model
func (e *Endpoint) SetMeta(d Meta) {
	e.Meta = d
}

// Validate validates the object
func (e *Endpoint) Validate() {
	e.AddError(ValidEndpointName("Invalid Id", e.Id))
}

// Prefix returns the type of object
func (e *Endpoint) Prefix() string {
	return "endpoints"
}

// Key returns the key for this object
func (e *Endpoint) Key() string {
	return e.Id
}

// KeyName returns the name of the field that is the key for this object
func (e *Endpoint) KeyName() string {
	return "Id"
}

// GetDescription returns the models Description
func (e *Endpoint) GetDescription() string {
	return e.Description
}

// Fill initials an Endpoint
func (e *Endpoint) Fill() {
	if e.Meta == nil {
		e.Meta = Meta{}
	}
	e.Validation.fill(e)
	if e.Params == nil {
		e.Params = map[string]interface{}{}
	}
	if e.VersionSets == nil {
		e.VersionSets = []string{}
	}
	if e.Plugins == nil {
		e.Plugins = []*Plugin{}
	}
	if e.Components == nil {
		e.Components = []*Element{}
	}
	if e.Prefs == nil {
		e.Prefs = map[string]string{}
	}
	if e.Global == nil {
		e.Global = map[string]interface{}{}
	}
	if e.Actions == nil {
		e.Actions = []*ElementAction{}
	}
}

// AuthKey returns the value of the key for auth purposes
func (e *Endpoint) AuthKey() string {
	return e.Key()
}

// SliceOf returns a slice of the model
func (e *Endpoint) SliceOf() interface{} {
	s := []*Endpoint{}
	return &s
}

// ToModels converts a slice of Endpoints into a slice of Model
func (e *Endpoint) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Endpoint)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

// GetParams returns the parameters on the Endpoint
// The returned map is a shallow copy.
func (e *Endpoint) GetParams() map[string]interface{} {
	return copyMap(e.Params)
}

// SetParams replaces the current parameters with a shallow
// copy of the input map.
func (e *Endpoint) SetParams(p map[string]interface{}) {
	e.Params = copyMap(p)
}

// CanHaveActions indicates that the model can have actions
func (e *Endpoint) CanHaveActions() bool {
	return true
}

// SetName sets the name. In this case, it sets Id.
func (e *Endpoint) SetName(name string) {
	e.Id = name
}
