package models

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	BsdpOS9       = 0
	BsdpOSX       = 1
	BsdpOSXServer = 2
	BsdpDiags     = 3
)

// BsdpBootOption contains the information needed to boot their
// systems using the aapl BDSP boot protocol.  Unless you are
// provisioning Mac systems, you probably don't need to care about
// this.
type BsdpBootOption struct {
	Index     uint16 `plist:"Index"`
	Install   bool   `plist:"IsInstall"`
	OSType    byte   `plist:"Kind"`
	OSVersion string `plist:"osVersion"`
	Name      string `plist:"Name"`
	Booter    string `plist:"BootFile"`
	RootPath  string `plist:"RootPath"`
}

func (bo *BsdpBootOption) OSName() string {
	switch bo.OSType {
	case BsdpOS9:
		return "os9"
	case BsdpOSX:
		return "osx"
	case BsdpOSXServer:
		return "osxsrv"
	case BsdpDiags:
		return "diags"
	default:
		return "unknown"
	}
}

func (bo *BsdpBootOption) InstallType() string {
	switch bo.Install {
	case true:
		return "netinstall"
	case false:
		return "netboot"
	}
	return "Impossible"
}

func (bo *BsdpBootOption) String() string {
	res := []string{
		bo.InstallType(),
		bo.OSName(),
		bo.OSVersion,
		fmt.Sprintf("%d", bo.Index),
		bo.Name,
		bo.Booter,
		bo.RootPath,
	}
	return strings.Join(res, ":")
}

func (bo *BsdpBootOption) MarshalText() ([]byte, error) {
	return []byte(bo.String()), nil
}

func (bo *BsdpBootOption) UnmarshalText(buf []byte) error {
	parts := strings.Split(string(buf), ":")
	if len(parts) != 7 {
		return fmt.Errorf("Want 7 parts, got %d", len(parts))
	}
	for i := range parts {
		val := parts[i]
		switch i {
		case 0:
			switch val {
			case "netboot":
				bo.Install = false
			case "netinstall":
				bo.Install = true
			default:
				return fmt.Errorf("Invalid install type %s", val)
			}
		case 1:
			switch val {
			case "os9":
				return fmt.Errorf("haha, try again. OS 9, really?")
			case "osx":
				bo.OSType = BsdpOSX
			case "osxsrv":
				bo.OSType = BsdpOSXServer
			case "diags":
				bo.OSType = BsdpDiags
			default:
				return fmt.Errorf("Unkown NBSP OS type %s", val)
			}
		case 2:
			bo.OSVersion = val
		case 3:
			idx, err := strconv.ParseInt(val, 0, 64)
			if err != nil {
				return err
			}
			bo.Index = uint16(idx)
		case 4:
			bo.Name = val
		case 5:
			bo.Booter = val
		case 6:
			bo.RootPath = val
		}
	}
	return nil
}
