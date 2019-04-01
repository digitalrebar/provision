// +build linux

package cli

import (
	"github.com/digitalrebar/provision/models"
	"github.com/rackn/gohai/plugins/dmi"
	"github.com/rackn/gohai/plugins/net"
	"github.com/rackn/gohai/plugins/storage"
	"github.com/rackn/gohai/plugins/system"
)

type dmiinfo interface {
	Class() string
}

func gohai() error {
	res := &models.Error{}
	infos := map[string]dmiinfo{}
	defer prettyPrint(infos)
	dmiInfo, err := dmi.Gather()
	if err != nil {
		res.Errorf("Failed to gather DMI information: %v", err)
	} else {
		infos[dmiInfo.Class()] = dmiInfo
	}
	netInfo, err := net.Gather()
	if err != nil {
		res.Errorf("Failed to gather network info: %v", err)
	} else {
		infos[netInfo.Class()] = netInfo
	}
	sysInfo, err := system.Gather()
	if err != nil {
		res.Errorf("Failed to gather basic OS info: %v", err)
	} else {
		infos[sysInfo.Class()] = sysInfo
	}
	storInfo, err := storage.Gather()
	if err != nil {
		res.Errorf("Failed to gather storage info: %v", err)
	} else {
		infos[storInfo.Class()] = storInfo
	}
	return res.HasError()
}
