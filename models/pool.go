package models

import (
	"bytes"
	"encoding/json"
	"log"
)

// PoolStatuses define the valid status for Machines in the pool
var PoolStatuses []string = []string{
	"Joining",
	"HoldJoin",
	"Free",
	"Building",
	"HoldBuild",
	"InUse",
	"Destroying",
	"HoldDestroy",
	"Leaving",
	"HoldLeave",
}

type PoolStatus string

const PS_JOINING = "Joining"
const PS_HOLD_JOIN = "HoldJoin"
const PS_FREE = "Free"
const PS_BUILDING = "Building"
const PS_HOLD_BUILD = "HoldBuild"
const PS_IN_USE = "InUse"
const PS_DESTROYING = "Destroying"
const PS_HOLD_DESTROY = "HoldDestroy"
const PS_LEAVING = "Leaving"
const PS_HOLD_LEAVE = "HoldLeave"

/*
 * Pool defines the basics of the pool.  This is a static object
 * that can be shared through content packs and version sets.
 *
 * Membership is dynamic and truth is from the machines' state.
 *
 * The transition actions are used on machines moving through the
 * pool (into and out of, allocated and released).
 *   EnterActions
 *   ExitActions
 *   AllocateActions
 *   ReleaseActions
 *
 * Params are used to provide default values.
 *
 * AutoFill Parameters are:
 *   UseAutoFill      bool
 *   MinFree          int32
 *   MaxFree          int32
 *   CreateParameters map[string]interface{}
 *   AcquirePool      string
 */
type Pool struct {
	Validation
	Access
	Meta
	Owned
	Bundled

	Id            string
	Description   string `json:",omitempty"`
	Documentation string `json:",omitempty"`
	ParentPool    string `json:",omitempty"`

	EnterActions    *PoolTransitionActions `json:",omitempty"`
	ExitActions     *PoolTransitionActions `json:",omitempty"`
	AllocateActions *PoolTransitionActions `json:",omitempty"`
	ReleaseActions  *PoolTransitionActions `json:",omitempty"`

	AutoFill *PoolAutoFill `json:",omitempty"`
}

// PoolTransitionActions define the default actions that should happen to a machine upon
// movement through the pool.
type PoolTransitionActions struct {
	Workflow         string                 `json:",omitempty"`
	AddProfiles      []string               `json:",omitempty"`
	AddParameters    map[string]interface{} `json:",omitempty"`
	RemoveProfiles   []string               `json:",omitempty"`
	RemoveParameters []string               `json:",omitempty"`
}

// PoolAutoFill are rules for dynamic pool sizing
type PoolAutoFill struct {
	UseAutoFill      bool                   `json:"UseAutoFill,omitempty"`
	MinFree          int32                  `json:"MinFree,omitempty"`
	MaxFree          int32                  `json:"MaxFree,omitempty"`
	CreateParameters map[string]interface{} `json:"CreateParameters,omitempty"`
	AcquirePool      string                 `json:"AcquirePool,omitempty"`
	ReturnPool       string                 `json:"ReturnPool,omitempty"`
}

// PoolResults is dynamically built provide membership and status.
type PoolResults map[PoolStatus][]*PoolResult

// PoolResult is the common return structure most operations
type PoolResult struct {
	Name      string
	Uuid      string
	Allocated bool
	Status    PoolStatus
}

func (p *Pool) Key() string {
	return p.Id
}

func (p *Pool) KeyName() string {
	return "Id"
}

func (p *Pool) AuthKey() string {
	return p.Key()
}

func (p *Pool) GetDescription() string {
	return p.Description
}

func (p *Pool) GetDocumentation() string {
	return p.Documentation
}

func (p *Pool) Prefix() string {
	return "pools"
}

// Clone the pool
func (p *Pool) Clone() *Pool {
	p2 := &Pool{}
	buf := bytes.Buffer{}
	enc, dec := json.NewEncoder(&buf), json.NewDecoder(&buf)
	if err := enc.Encode(p); err != nil {
		log.Panicf("Failed to encode pools:%s: %v", p.Id, err)
	}
	if err := dec.Decode(p2); err != nil {
		log.Panicf("Failed to decode pools:%s: %v", p.Id, err)
	}
	return p2
}

func (p *Pool) Fill() {
	if p.Meta == nil {
		p.Meta = Meta{}
	}
	if p.Errors == nil {
		p.Errors = []string{}
	}
}

func (p *Pool) CanHaveActions() bool {
	return true
}

func (p *Pool) SliceOf() interface{} {
	s := []*Pool{}
	return &s
}

func (p *Pool) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Pool)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

// SetName sets the name. In this case, it sets Id.
func (p *Pool) SetName(name string) {
	p.Id = name
}
