package api

import (
	"net"
	"runtime"
	"testing"

	"github.com/digitalrebar/provision"
	"github.com/digitalrebar/provision/models"
)

func TestInfo(t *testing.T) {
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
			Stats: []*models.Stat{
				{
					Name:  "machines.count",
					Count: 0,
				},
				{
					Name:  "subnets.count",
					Count: 0,
				},
			},
			Arch:    runtime.GOARCH,
			Os:      runtime.GOOS,
			Version: provision.RS_VERSION,
			Id:      "Fred",
			Features: []string{
				"api-v3",
				"sane-exit-codes",
				"common-blob-size",
				"change-stage-map",
				"job-exit-states",
				"package-repository-handling",
				"profileless-machine",
				"threaded-log-levels",
			},
		},
		expectErr: nil,
		op: func() (interface{}, error) {
			return session.Info()
		},
	}
	test.run(t)

}
