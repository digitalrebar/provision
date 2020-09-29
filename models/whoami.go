package models

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/rackn/gohai/plugins/dmi"
	ghnet "github.com/rackn/gohai/plugins/net"

	"github.com/pborman/uuid"
)

type WhoamiResult struct {
	Score int       `json:",omitempty"`
	Uuid  uuid.UUID `json:",omitempty"`
	Token string    `json:",omitempty"`
}

// Whoami contains the elements used toi fingerprint a machine, along with
// the results of the fingerprint comparison request
type Whoami struct {
	Result      WhoamiResult
	Fingerprint MachineFingerprint
	MacAddrs    []string
	OnDiskUUID  string
}

// From Ubuntu's fwts package.  We should update these every once in a while.

var badSerials = []string{
	"0000000",
	"00000000",
	"000000000",
	"0000000000",
	"012345678",
	"0123456789",
	"01234567890",
	"012345678900",
	"0123456789000",
	"0x00000000",
	"0x0000000000000000",
	"<cut out>",
	"Base Board Serial Number",
	"Chassis Serial Number",
	"Empty",
	"MB-1234567890",
	"N/A",
	"NA",
	"NB-0123456789",
	"NB-1234567890",
	"None",
	"None1",
	"Not Available",
	"Not Specified",
	"Not Supported by CPU",
	"Not Supported",
	"NotSupport",
	"OEM Chassis Serial Number",
	"OEM_Define1",
	"OEM_Define2",
	"OEM_Define3",
	"SerNum0",
	"SerNum00",
	"SerNum01",
	"SerNum02",
	"SerNum03",
	"SerNum1",
	"SerNum2",
	"SerNum3",
	"SerNum4",
	"System Serial Number",
	"TBD by ODM",
	"To Be Defined By O.E.M",
	"To Be Filled By O.E.M.",
	"To be filled by O.E.M.",
	"Unknow",
	"Unknown",
	"XXXXX",
	"XXXXXX",
	"XXXXXXXX",
	"XXXXXXXXX",
	"XXXXXXXXXX",
	"XXXXXXXXXXX",
	"XXXXXXXXXXXX",
	"[Empty]",
}

var badAssets = []string{
	"0000000000",
	"0x00000000",
	"1234567890",
	"123456789000",
	"9876543210",
	"<cut out>",
	"A1_AssetTagNum0",
	"A1_AssetTagNum1",
	"A1_AssetTagNum2",
	"A1_AssetTagNum3",
	"ABCDEFGHIJKLM",
	"ATN12345678901234567",
	"Asset Tag",
	"Asset Tag:",
	"Asset tracking",
	"Asset-1234567890",
	"AssetTagNum0",
	"AssetTagNum1",
	"AssetTagNum2",
	"AssetTagNum3",
	"AssetTagNum4",
	"Base Board Asset Tag",
	"Base Board Asset Tag#",
	"Chassis Asset Tag",
	"Fill By OEM",
	"N/A",
	"No Asset Information",
	"No Asset Tag",
	"None",
	"Not Available",
	"Not Specified",
	"OEM_Define0",
	"OEM_Define1",
	"OEM_Define2",
	"OEM_Define3",
	"OEM_Define4",
	"TBD by ODM",
	"To Be Defined By O.E.M",
	"To Be Filled By O.E.M.",
	"To be filled by O.E.M.",
	"Unknown",
	"XXXXXX",
	"XXXXXXX",
	"XXXXXXXX",
	"XXXXXXXXX",
	"XXXXXXXXXX",
	"XXXXXXXXXXX",
	"XXXXXXXXXXXX",
}

func vIsUseable(v string, bads []string) bool {
	vv := strings.TrimSpace(v)
	i := sort.SearchStrings(bads, vv)
	return len(vv) > 0 && !(i < len(bads) && bads[i] == vv)
}

func vIsSane(v string, bads []string) bool {
	vv := strings.TrimSpace(v)
	i := sort.SearchStrings(bads, vv)
	return !(i < len(bads) && bads[i] == vv)
}

func uuidIsSane(v string) bool {
	if strings.HasSuffix(v, "-000000000000") {
		// suffix of all zeros?  No bueno.
		return false
	}
	if v == "0A0A0A0A-0A0A-0A0A-0A0A-0A0A0A0A0A0A" {
		// Specifically a sentinel value.
		return false
	}
	id := uuid.Parse(v)
	if id == nil {
		return false
	}
	_, ok := id.Version()
	return ok
}

// Fill fills out the MacAddrs list and the MachineFingerprint from
// the DMI information on the machine.
func (w *Whoami) Fill() error {
	(&w.Fingerprint).Fill()
	w.MacAddrs = []string{}
	hasher := sha256.New()
	dmiinfo, err := dmi.Gather()
	if err != nil {
		return err
	}
	if dmiinfo.Hypervisor == "" {
		// we can only trust that this information will be unique and independent from
		// another if running on physical hardware.  This may be a false positive as
		// new hypervisors are released, but it detects all the common ones and several
		// uncommon ones.  Too bad that it means we have to fall back to relying on MAC address
		// based uniqueness checking, but that is pretty much where we were earlier.
		if vIsUseable(dmiinfo.System.SerialNumber, badSerials) {
			fmt.Fprint(hasher, dmiinfo.System.Manufacturer, dmiinfo.System.ProductName, dmiinfo.System.SerialNumber)
			w.Fingerprint.SSNHash = hasher.Sum(nil)
			hasher.Reset()
		}
		if len(dmiinfo.Chassis) == 1 && vIsUseable(dmiinfo.Chassis[0].SerialNumber, badSerials) {
			fmt.Fprint(hasher, dmiinfo.Chassis[0].Manufacturer, dmiinfo.Chassis[0].SerialNumber)
			w.Fingerprint.CSNHash = hasher.Sum(nil)
			hasher.Reset()
		} else if len(dmiinfo.Baseboards) == 1 && vIsUseable(dmiinfo.Baseboards[0].SerialNumber, badSerials) {
			fmt.Fprint(hasher, dmiinfo.Baseboards[0].Manufacturer, dmiinfo.Baseboards[0].ProductName, dmiinfo.Baseboards[0].SerialNumber)
			w.Fingerprint.CSNHash = hasher.Sum(nil)
			hasher.Reset()
		}
		if uuidIsSane(dmiinfo.System.UUID) {
			w.Fingerprint.SystemUUID = dmiinfo.System.UUID
		}
		for _, mem := range dmiinfo.Memory.Devices {
			// For now, we just assume that if the dimm has an invalid asset tag then we
			// should also assume that the serial number is not going to be unique.
			// This holds true for Corsair memory, at least.
			if !(vIsUseable(mem.SerialNumber, badSerials) && vIsSane(mem.AssetTag, badAssets)) {
				continue
			}
			fmt.Fprint(hasher, mem.Manufacturer, mem.PartNumber, mem.SerialNumber)
			w.Fingerprint.MemoryIds = append(w.Fingerprint.MemoryIds, hasher.Sum(nil))
			hasher.Reset()
		}
	}
	sort.Slice(w.Fingerprint.MemoryIds, func(i, j int) bool {
		return bytes.Compare(w.Fingerprint.MemoryIds[i], w.Fingerprint.MemoryIds[j]) == -1
	})
	if len(w.Fingerprint.MemoryIds) > 1 {
		for idx, v := range w.Fingerprint.MemoryIds[1:] {
			if bytes.Equal(v, w.Fingerprint.MemoryIds[idx]) {
				// Sigh, we also cannot trust that the serial number is really a serial number.
				w.Fingerprint.MemoryIds = [][]byte{}
				break
			}
		}
	}
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

// FromMachine extracts the Fingerprint and HardwareAddrs fields from a Machine
// and populates Whoami with it.
func (w *Whoami) FromMachine(m *Machine) {
	w.MacAddrs = append([]string{}, m.HardwareAddrs...)
	w.Fingerprint = m.Fingerprint
}

// ToMachine saves the Fingerprint and the MacAddrs fields onto the passed Machine.
func (w *Whoami) ToMachine(m *Machine) {
	m.HardwareAddrs = append([]string{}, w.MacAddrs...)
	m.Fingerprint = w.Fingerprint
}

// Score calculates how closely the passed in Whoami matches a candidate Machine.
// In the current implementation, Score awards points based on the following
// criteria:
//
// * 25 points if the Machine has an SSNHash that matches the one in the Whoami
//
// * 25 points if the Machine has a CSNHash that matches the one in the Whoami
//
// * 50 points if the Machine has a SystemUUID that matches the one in the Whoami
//
// * 0 to 100 points varying depending on how many memory DIMMs from the machine
//   fingerprint are present in Whoami.
//
// * 0 to 100 points varying depending on how many HardwareAddrs from the Machine
//   are present in Whoami.
//
// * 1000 points if the machine UUID matches OnDiskUUID
//
// If the score is less than 100 at the end of the scoring process, it is rounded down
// to zero.  The intent is to be resilient in the face of hardware changes:
//
// * SSNHash, CSNHash, and SystemUUID come from the motherboard.
//
// * MemoryIds are generated deterministically from the DIMMs installed in the system
//
// * MacAddrs comes from the physical Ethernet devices in the system.
func (w *Whoami) Score(m *Machine) (score int) {
	if len(m.Fingerprint.SSNHash) > 0 {
		if bytes.Equal(w.Fingerprint.SSNHash, m.Fingerprint.SSNHash) {
			score += 25
		}
	}
	if len(m.Fingerprint.CSNHash) > 0 {
		if bytes.Equal(w.Fingerprint.CSNHash, m.Fingerprint.CSNHash) {
			score += 25
		}
	}
	if m.Fingerprint.SystemUUID != "" {
		if w.Fingerprint.SystemUUID == m.Fingerprint.SystemUUID {
			score += 50
		}
	}
	var matched int
	var j int
	if len(m.Fingerprint.MemoryIds) > 0 {
		for _, probe := range m.Fingerprint.MemoryIds {
			for j = range w.Fingerprint.MemoryIds {
				cmp := bytes.Compare(w.Fingerprint.MemoryIds[j], probe)
				if cmp == 0 {
					matched++
					break
				}
			}
		}
		score += (100 * matched) / len(m.Fingerprint.MemoryIds)
	}
	if len(m.HardwareAddrs) > 0 {
		matched = 0
		for _, probe := range m.HardwareAddrs {
			for j = range w.MacAddrs {
				cmp := strings.Compare(w.MacAddrs[j], probe)
				if cmp == 0 {
					matched++
					break
				}
			}
		}
		score += (100 * matched) / len(m.HardwareAddrs)
	}
	if m.UUID() == w.OnDiskUUID {
		score += 1000
	}
	if score < 100 {
		score = 0
	}
	return
}
