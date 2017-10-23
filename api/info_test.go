package api

import (
	"runtime"
	"testing"

	"github.com/digitalrebar/provision"
	"github.com/digitalrebar/provision/models"
)

func TestInfo(t *testing.T) {
	test := &crudTest{
		name: "get info",
		expectRes: &models.Info{
			ApiPort:            10001,
			FilePort:           10002,
			ProvisionerEnabled: true,
			TftpEnabled:        true,
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
			Arch:     runtime.GOARCH,
			Os:       runtime.GOOS,
			Version:  provision.RS_VERSION,
			Id:       "Fred",
			Features: []string{"api-v3", "sane-exit-codes", "common-blob-size"},
		},
		expectErr: nil,
		op: func() (interface{}, error) {
			return session.Info()
		},
	}
	test.run(t)

}
