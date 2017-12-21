package models

import "net"

// swagger:model
type Stat struct {
	// required: true
	Name string `json:"name"`
	// required: true
	Count int `json:"count"`
}

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
	Stats    []*Stat  `json:"stats"`
	Features []string `json:"features"`
}
