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
	// required: true
	Electable bool `json:"electable"`
}

// Stat contains a basic statistic about dr-provision
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
	// Arch is the system architecture of the running dr-provision endpoint.
	// It is the same value that would be return by runtime.GOARCH
	// required: true
	Arch string `json:"arch"`
	// Os is the operating system the dr-provision endpoint is running on.
	// It is the same value returned by runtime.GOARCH
	// required: true
	Os string `json:"os"`
	// Version is the full version of dr-provision.
	// required: true
	Version string `json:"version"`
	// Id is the local ID for this dr-provision.  If not overridden by
	// an environment variable or a command line argument, it will
	// be the lowest MAC address of all the physical nics attached to the system.
	// required: true
	Id string `json:"id"`
	// LocalId is the same as Id, except it is always the MAC address form.
	// required: true
	LocalId string `json:"local_id"`
	// HaId is the user-assigned high-availability ID for this endpoint.
	// All endpoints in the same HA cluster must have the same HaId.
	// required: true
	HaId string `json:"ha_id"`
	// ApiPort is the TCP port that the API lives on.  Defaults to 8092
	// required: true
	ApiPort int `json:"api_port"`
	// FilePort is the TCP port that the static file HTTP server lives on.
	// Defaults to 8091
	// required: true
	FilePort int `json:"file_port"`
	// SecureFilePort is the TCP port that the static file HTTPS server lives on.
	// Defaults to 8090
	// required: true
	SecureFilePort int `json:"secure_file_port"`
	// DhcpPort is the UDP port that the DHCPv4 server listens on.
	// Defaults to 67
	// required: true
	DhcpPort int `json:"dhcp_port"`
	// BinlPort is the UDP port that the BINL server listens on.
	// Defaults to 4011
	// required: true
	BinlPort int `json:"binl_port"`
	// TftpPort is the UDP port that the TFTP server listens on.
	// Defaults to 69, dude.
	// required: true
	TftpPort int `json:"tftp_port"`
	// TftpEnabled is true if the TFTP server is enabled.
	// required: true
	TftpEnabled bool `json:"tftp_enabled"`
	// DhcpEnabled is true if the DHCP server is enabled.
	// required: true
	DhcpEnabled bool `json:"dhcp_enabled"`
	// BinlEnabled is true if the BINL server is enabled.
	// required: true
	BinlEnabled bool `json:"binl_enabled"`
	// ProvisionerEnabled is true if the static file HTTP server is enabled.
	// required: true
	ProvisionerEnabled bool `json:"prov_enabled"`
	// SecureProvisionerEnabled is true if the static file HTTPS server is enabled.
	// required: true
	SecureProvisionerEnabled bool `json:"secure_prov_enabled"`
	// Address is the IP address that the system appears to listen on.
	// If a default address was assigned via environment variable or command line,
	// it will be that address, otherwise it will be the IP address of the interface
	// that has the default IPv4 route.
	// required: true
	Address net.IP `json:"address"`
	// Manager indicates whether this dr-provision can act as a manager of
	// other dr-provision instances.
	// required: true
	Manager bool `json:"manager"`
	// Stats lists some basic object statistics.
	// required: true
	Stats []Stat `json:"stats"`
	// Features is a list of features implemented in this dr-provision endpoint.
	// Clients should use this field when determining what features are available
	// on anu given dr-provision instance.
	Features []string `json:"features"`
	// Scopes lists all static permission scopes available.
	Scopes map[string]map[string]struct{} `json:"scopes"`
	// License is an embedded copy of the licenses present on the system.
	License LicenseBundle

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

	// ClusterState is the current state of the consensus cluster that this
	// node is a member of.  As of v4.6, all nodes are in at least a single-node cluster
	ClusterState ClusterState
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
