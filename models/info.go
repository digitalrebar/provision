package models

import (
	"net"
	"runtime"
)

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
	Stats    []Stat                         `json:"stats"`
	Features []string                       `json:"features"`
	Scopes   map[string]map[string]struct{} `json:"scopes"`
	License  LicenseBundle
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
		i.Features = []string{
			"api-v3",
			"sane-exit-codes",
			"common-blob-size",
			"change-stage-map",
			"job-exit-states",
			"package-repository-handling",
			"profileless-machine",
			"threaded-log-levels",
			"plugin-v2",
			"fsm-runner",
			"plugin-v2-safe-config",
			"workflows",
			"default-workflow",
			"http-range-header",
			"roles",
			"tenants",
			"secure-params",
			"separate-meta-api",
			"slim-objects",
			"secure-param-upgrade",
			"sprig",
			"multiarch",
			"actions-in-task-list",
			"endpoint-refs",
			"endpoint-proxy",
			"inline-upgrade",
			"bundle-objects",
			"secure-params-in-content-packs",
			"task-prerequisites",
			"content-prerequisite-version-checking",
			"stage-paramer",
			"auto-boot-target",
			"partial-objects",
			"regex-string-filters",
			"file-iso-exists-info-render",
			"log-has-head-method",
		}
	}
	if i.Scopes == nil {
		scopes := map[string]map[string]struct{}{}
		actionScopeLock.Lock()
		defer actionScopeLock.Unlock()
		Remarshal(allScopes, &scopes)
		i.Scopes = scopes
	}
}
