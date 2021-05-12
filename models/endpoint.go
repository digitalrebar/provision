package models

import (
	"fmt"
)

/*
 * Endpoint tracks who to get events from and where to send them when
 * actions are required.
 */

// Element define a part of the endpoint
// This can be a file, a pref, global profile parameter,
// DRP itself, content packages, plugin_providers, or plugins.
type Element struct {
	// Type defines the type of element
	// This can be:
	//   DRP, DRPUX, File, Global, Plugin, Pref, PluginProvider, ContentPackage
	Type string
	// Version defines the short or reference version of the element.
	// e.g. tip, stable, v4.3.6
	Version string
	// Name defines the name of the element.  Normally, this is the
	// name of the the DRP, DRPUX, filename, plugin, ContentPackage, or PluginProvider Name.
	// For Global and Pref, these are the name of the global parameter or preference.
	Name string
	// ActualVersion is the actual catalog version referenced by this element.
	// This is used for translating tip and stable into a real version.
	// This is the source of the file element.  This can be a relative or absolute path or an URL.
	ActualVersion string
}

// ElementAction defines an action to take on an Element
type ElementAction struct {
	Element
	// Action defines what is to be done to this element.
	// These can be Set for Pref, Global.
	// These can be AddOrUpdate and Delete for the reset of the elements.
	// This field is ignored for the DRP and DRPUX element.  It is assumed AddOrUpdate.
	Action string
	// Value defines what should be set or applied.  This field is used
	// for the  plugin, pref, global, and file elements.
	//
	// Plugin, Pref, and Global elements use this as the value for the element.
	// File elements use this field to determine if it should be exploded.
	Value interface{}
}

// String prints a user-friendly format of an ElementAction
func (ea *ElementAction) String() string {
	return fmt.Sprintf("%s %s %s:%s", ea.Action, ea.Type, ea.Name, ea.Version)
}

// Endpoint represents a managed Endpoint
//
// This object is used to reflect the current state of a downstream endpoint.
//
// It also shows the desired configuration state of the downstream endpoint
// through the applied versions sets.
//
// It acts as the control point for applying updates through the Apply field.
// A user can also "dry-run" a set of changes to see what would happen by viewing
// the Actions field.  This also shows remaining work while Apply is set to true.
//
// It is similar to a machine object in that it has parameters that define
// how to access the endpoint and its state.
//
// swagger:model
type Endpoint struct {
	Validation
	Access
	Meta
	Owned
	Bundled

	// Id is the name of the DRP endpoint this should match the HA pair's ID or the DRP ID of a single node.
	Id string `index:",key"`

	// Description is a string for providing a simple description
	Description string `json:"Description,omitempty"`
	// Documentation is a string for providing additional in depth information.
	Documentation string `json:"Documentation,omitempty"`

	// Params holds the access parameters - these should be secure parameters.
	// They are:
	//   manager/username
	//   manager/password
	//   manager/url
	Params map[string]interface{} `json:"Params,omitempty"`

	// ConnectionStatus reflects the manager's state of interaction with the endpoint
	ConnectionStatus string `json:"ConnectionStatus,omitempty"`

	// VersionSet - Deprecated - was a single version set.
	// This should be specified within the VersionSets list
	VersionSet string `json:"VersionSet,omitempty"`

	// VersionSets replaces VersionSet - code processes both
	// This is the list of version sets to apply.  These are merged
	// with the first in the list having priority over later elements in the list.
	VersionSets []string `json:"VersionSets,omitempty"`

	// Apply toggles whether the manager should update the endpoint.
	Apply bool `json:"Apply,omitempty"`
	// HaId is the HaId of the endpoint
	HaId string `json:"HaId,omitempty"`
	// Arch is the arch of the endpoint - Golang arch format.
	Arch string `json:"Arch,omitempty"`
	// Os is the os of the endpoint - Golang os format.
	Os string `json:"Os,omitempty"`
	// DRPVersion is the version of the drp endpoint running.
	DRPVersion string `json:"DRPVersion,omitempty"`
	// DRPUXVersion is the version of the ux installed on the endpoint.
	DRPUXVersion string `json:"DRPUXVersion,omitempty"`
	// Components is the list of ContentPackages and PluginProviders installed
	// and their versions
	Components []*Element `json:"Components,omitempty"`
	// Plugins is the list of Plugins configured on the endpoint.
	Plugins []*Plugin `json:"Plugins,omitempty"`
	// Prefs is the value of all the prefs on the endpoint.
	Prefs map[string]string `json:"Prefs,omitempty"`
	// Global is the Parameters of the global profile.
	Global map[string]interface{} `json:"Global,omitempty"`
	// Actions is the list of actions to take to make the endpoint
	// match the version sets on in the endpoint object.
	Actions []*ElementAction `json:"Actions,omitempty"`
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
