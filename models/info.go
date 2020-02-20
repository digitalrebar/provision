package models

import (
	"net"
	"runtime"
)

// HaPassiveState the state of the passive node
//
// swagger:model
type HaPassiveState struct {
	// required: true
	Id string `json:"id"`
	// required: true
	Address string `json:"address"`
	// required: true
	State string `json:"state"`
}

// Stat contains a basic statistic sbout dr-provision
//
// swagger:model
type Stat struct {
	// required: true
	Name string `json:"name"`
	// required: true
	Count int `json:"count"`
}

// Info contains information on how the running instance of
// dr-provision is configured.
//
// For passive nodes, the license, scopes, and stats are not filled in.
//
// swagger:model
type Info struct {
	// required: true
	Arch string `json:"arch"`
	// required: true
	Os string `json:"os"`
	// required: true
	Version string `json:"version"`
	// required: true
	Id string `json:"id"`
	// required: true
	LocalId string `json:"local_id"`
	// required: true
	HaId string `json:"ha_id"`
	// required: true
	ApiPort int `json:"api_port"`
	// required: true
	FilePort int `json:"file_port"`
	// required: true
	DhcpPort int `json:"dhcp_port"`
	// required: true
	BinlPort int `json:"binl_port"`
	// required: true
	TftpPort int `json:"tftp_port"`
	// required: true
	TftpEnabled bool `json:"tftp_enabled"`
	// required: true
	DhcpEnabled bool `json:"dhcp_enabled"`
	// required: true
	BinlEnabled bool `json:"binl_enabled"`
	// required: true
	ProvisionerEnabled bool `json:"prov_enabled"`
	// required: true
	Address net.IP `json:"address"`
	// required: true
	Manager bool `json:"manager"`
	// required: true
	Stats    []Stat                         `json:"stats"`
	Features []string                       `json:"features"`
	Scopes   map[string]map[string]struct{} `json:"scopes"`
	License  LicenseBundle

	// Errors returns the current system errors.
	// required: true
	Errors []string `json:"errors"`

	// HaEnabled indicates if High Availability is enabled
	HaEnabled bool `json:"ha_enabled"`
	// HaVirtualAddress is the Virtual IP Address of the systems
	HaVirtualAddress string `json:"ha_virtual_address"`
	// HaIsActive indicates Active (true) or Passive (false)
	// required: true
	HaIsActive bool `json:"ha_is_active"`
	// HaStatus indicates current state
	// For Active, Up is the only value.
	// For Passive, Connecting, Syncing, In-Sync
	// required: true
	HaStatus string `json:"ha_status"`

	// HaActiveId is the id of current active node
	HaActiveId string `json:"ha_active_id"`
	// HaPassiveState is a list of passive node's and their current state
	// This is only valid from the Active node
	HaPassiveState []*HaPassiveState `json:"ha_passive_state"`
}

// HasFeature is a helper function to determine if a requested feature
// is present.
func (i *Info) HasFeature(f string) bool {
	for _, v := range i.Features {
		if v == f {
			return true
		}
	}
	return false
}

func (i *Info) Fill() {
	i.Arch = runtime.GOARCH
	i.Os = runtime.GOOS
	if i.Stats == nil {
		i.Stats = make([]Stat, 0, 0)
	}
	if i.Features == nil {
		i.Features = []string{}
	}
	if i.Scopes == nil {
		scopes := map[string]map[string]struct{}{}
		actionScopeLock.Lock()
		defer actionScopeLock.Unlock()
		Remarshal(allScopes, &scopes)
		i.Scopes = scopes
	}
	if i.Errors == nil {
		i.Errors = []string{}
	}
	if i.HaPassiveState == nil {
		i.HaPassiveState = []*HaPassiveState{}
	}
}

func (i *Info) AddUpdatePassive(id, address, state string) {
	if i.HaPassiveState == nil {
		i.HaPassiveState = []*HaPassiveState{&HaPassiveState{Id: id, Address: address, State: state}}
		return
	}
	for _, ps := range i.HaPassiveState {
		if ps.Id == id {
			ps.State = state
			if address != "" {
				ps.Address = address
			}
			return
		}
	}
	i.HaPassiveState = append(i.HaPassiveState, &HaPassiveState{Id: id, Address: address, State: state})
}

func (i *Info) RemovePassive(id string) {
	if i.HaPassiveState == nil {
		return
	}
	idx := -1
	for ii, ps := range i.HaPassiveState {
		if ps.Id == id {
			idx = ii
		}
	}
	if idx != -1 {
		i.HaPassiveState[idx] = i.HaPassiveState[len(i.HaPassiveState)-1]
		i.HaPassiveState[len(i.HaPassiveState)-1] = nil
		i.HaPassiveState = i.HaPassiveState[:len(i.HaPassiveState)-1]
	}
}
