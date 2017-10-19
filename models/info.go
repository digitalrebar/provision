package models

// swagger:model
type Stat struct {
	// required: true
	Name string
	// required: true
	Count int
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
	TftpEnabled bool `json:"tftp_enabled"`
	// required: true
	DhcpEnabled bool `json:"dhcp_enabled"`
	// required: true
	ProvisionerEnabled bool `json:"prov_enabled"`
	// required: true
	Stats    []*Stat `json:"stats"`
	Features []string
}
