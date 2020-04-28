package models

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/rackn/gohai/plugins/dmi"
	ghnet "github.com/rackn/gohai/plugins/net"

	"github.com/google/uuid"
)

type Whoami struct {
	Result struct {
		Score int       `json:",omitempty"`
		Uuid  uuid.UUID `json:",omitempty"`
		Token string    `json:",omitempty"`
	}
	Fingerprint MachineFingerprint
	MacAddrs    []string
	OnDiskUUID  string
}

func (w *Whoami) Fill() error {
	(&w.Fingerprint).Fill()
	w.MacAddrs = []string{}
	hasher := sha256.New()
	dmiinfo, err := dmi.Gather()
	if err != nil {
		return err
	}
	if dmiinfo.System.SerialNumber != "" {
		fmt.Fprint(hasher, dmiinfo.System.Manufacturer, dmiinfo.System.ProductName, dmiinfo.System.SerialNumber)
		w.Fingerprint.SSNHash = hasher.Sum(nil)
		hasher.Reset()
	}
	if len(dmiinfo.Chassis) > 0 && dmiinfo.Chassis[0].SerialNumber != "" {
		fmt.Fprint(hasher, dmiinfo.System.Manufacturer, dmiinfo.System.ProductName, dmiinfo.Chassis[0].SerialNumber)
		w.Fingerprint.CSNHash = hasher.Sum(nil)
		hasher.Reset()
	}
	w.Fingerprint.SystemUUID = dmiinfo.System.UUID
	for _, mem := range dmiinfo.Memory.Devices {
		if mem.SerialNumber == "" {
			continue
		}
		fmt.Fprint(hasher, mem.Manufacturer, mem.PartNumber, mem.SerialNumber)
		w.Fingerprint.MemoryIds = append(w.Fingerprint.MemoryIds, hasher.Sum(nil))
		hasher.Reset()
	}
	sort.Slice(w.Fingerprint.MemoryIds, func(i, j int) bool {
		return bytes.Compare(w.Fingerprint.MemoryIds[i], w.Fingerprint.MemoryIds[j]) == -1
	})
	netinfo, err := ghnet.Gather()
	if err != nil {
		return err
	}
	for _, intf := range netinfo.Interfaces {
		if !intf.Sys.IsPhysical {
			continue
		}
		w.MacAddrs = append(w.MacAddrs, intf.HardwareAddr.String())
	}
	sort.Strings(w.MacAddrs)
	return nil
}

func (w *Whoami) FromMachine(m *Machine) {
	w.MacAddrs = append([]string{}, m.HardwareAddrs...)
	w.Fingerprint = m.Fingerprint
}

func (w *Whoami) ToMachine(m *Machine) {
	m.HardwareAddrs = append([]string{}, w.MacAddrs...)
	m.Fingerprint = w.Fingerprint
}

func (w *Whoami) Score(m *Machine) int {
	res := 0
	if len(w.Fingerprint.SSNHash) > 0 && bytes.Equal(w.Fingerprint.SSNHash, m.Fingerprint.SSNHash) {
		res += 25
	}
	if len(w.Fingerprint.CSNHash) > 0 && bytes.Equal(w.Fingerprint.CSNHash, m.Fingerprint.CSNHash) {
		res += 25
	}

	if w.Fingerprint.SystemUUID != "" && w.Fingerprint.SystemUUID == m.Fingerprint.SystemUUID {
		res += 50
	}
	var matched int
	var j int
	for _, probe := range m.Fingerprint.MemoryIds {
		for j = range w.Fingerprint.MemoryIds {
			cmp := bytes.Compare(w.Fingerprint.MemoryIds[j], probe)
			if cmp == 0 {
				matched++
				break
			}
			if cmp == 1 {
				break
			}
		}
	}
	res += int((float32(matched) / float32(len(m.Fingerprint.MemoryIds))) * 100)
	matched = 0
	for _, probe := range m.HardwareAddrs {
		for j = range w.MacAddrs {
			cmp := strings.Compare(w.MacAddrs[j], probe)
			if cmp == 0 {
				matched++
				break
			}
			if cmp == 1 {
				break
			}
		}
	}
	res += int((float32(matched) / float32(len(m.HardwareAddrs))) * 100)
	if m.UUID() == w.OnDiskUUID {
		res += 1000
	}
	if res < 100 {
		return 0
	}
	return res
}
