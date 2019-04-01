// +build linux

package cli

import (
	"fmt"

	"github.com/rackn/gohai/plugins/dmi"
	"github.com/rackn/gohai/plugins/net"
	"github.com/rackn/gohai/plugins/storage"
	"github.com/rackn/gohai/plugins/system"
)

type dmiinfo interface {
	Class() string
}

func gohai() error {
	infos := map[string]dmiinfo{}
	defer prettyPrint(infos)
	dmiInfo, err := dmi.Gather()
	if err != nil {
		return fmt.Errorf("Failed to gather DMI information: %v", err)
	}
	infos[dmiInfo.Class()] = dmiInfo
	netInfo, err := net.Gather()
	if err != nil {
		return fmt.Errorf("Failed to gather network info: %v", err)
	}
	infos[netInfo.Class()] = netInfo
	sysInfo, err := system.Gather()
	if err != nil {
		return fmt.Errorf("Failed to gather basic OS info: %v", err)
	}
	infos[sysInfo.Class()] = sysInfo
	storInfo, err := storage.Gather()
	if err != nil {
		return fmt.Errorf("Failed tp gather storage info: %v", err)
	}
	infos[storInfo.Class()] = storInfo
	return nil
}
