package api

import (
	"net"
	"runtime"
	"strings"
	"testing"

	"github.com/digitalrebar/provision/v4/models"
)

func TestInfo(t *testing.T) {
	localId := ""
	intfs, _ := net.Interfaces()
	for _, intf := range intfs {
		if (intf.Flags & net.FlagLoopback) == net.FlagLoopback {
			continue
		}
		if (intf.Flags & net.FlagUp) != net.FlagUp {
			continue
		}
		if strings.HasPrefix(intf.Name, "veth") {
			continue
		}
		localId = intf.HardwareAddr.String()
		break
	}

	test := &crudTest{
		name: "get info",
		expectRes: &models.Info{
			Address:            net.IPv4(127, 0, 0, 1),
			ApiPort:            10011,
			FilePort:           10012,
			BinlPort:           10015,
			DhcpPort:           10014,
			TftpPort:           10013,
			ProvisionerEnabled: true,
			TftpEnabled:        true,
			BinlEnabled:        true,
			DhcpEnabled:        true,
			Stats: []models.Stat{
				{
					Name:  "machines.count",
					Count: 0,
				},
				{
					Name:  "subnets.count",
					Count: 0,
				},
				{
					Name:  "contexts.count",
					Count: 0,
				},
			},
			Arch:           runtime.GOARCH,
			Os:             runtime.GOOS,
			Version:        "",
			HaId:           "Fred",
			Id:             "Fred",
			LocalId:        localId,
			Features:       []string{},
			License:        models.LicenseBundle{Licenses: []models.License{}},
			Scopes:         map[string]map[string]struct{}{},
			HaActiveId:     "Fred",
			HaEnabled:      false,
			HaIsActive:     true,
			HaPassiveState: []*models.HaPassiveState{},
			HaStatus:       "Up",
			Errors:         []string{},
		},
		expectErr: nil,
		op: func() (interface{}, error) {
			info, err := session.Info()
			if info != nil {
				info.Version = ""
				info.License = models.LicenseBundle{Licenses: []models.License{}}
				info.Features = []string{}
				info.Scopes = map[string]map[string]struct{}{}
				info.ClusterState = models.ClusterState{}
			}
			return info, err
		},
	}
	test.run(t)

}
