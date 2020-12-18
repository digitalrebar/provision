package models

import (
	"net"
	"strconv"
	"strings"
)

// ZoneRecord contains an individual record for a DNS zone
//
// swagger:model
type ZoneRecord struct {
	Type  uint16   // Type of record
	Name  string   // Name of record - This can contain template pieces
	Value []string // Value of record - This can contain template pieces.
}

// ZoneFilter contains a filter for the DNS packet to apply to the zone.
//
// swagger:model
type ZoneFilter struct {
	Type   uint16 // 1 = subnet, ...
	Filter string // filter string dependent upon Type
}

// Zone contains a list of records for DNS.
//
// swagger:model
type Zone struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	Name          string
	Description   string
	Documentation string

	Zone       string        // Base zone for this zone
	Filters    []*ZoneFilter // Matchers
	Priority   int           // Lower first
	Forwarders []string      // Forwarders for this zone.
	Continue   bool          // Should the next zone be processed.

	Records []*ZoneRecord
}

func (z *Zone) GetMeta() Meta {
	return z.Meta
}

func (z *Zone) SetMeta(d Meta) {
	z.Meta = d
}

// GetDocumentaiton returns the object's Documentation
func (z *Zone) GetDocumentation() string {
	return z.Documentation
}

// GetDescription returns the object's Description
func (z *Zone) GetDescription() string {
	return z.Description
}

func (z *Zone) Prefix() string {
	return "zones"
}

func (z *Zone) Key() string {
	return z.Name
}

func (z *Zone) KeyName() string {
	return "Name"
}

func (z *Zone) Fill() {
	z.Validation.fill(z)
	if z.Meta == nil {
		z.Meta = Meta{}
	}
	if z.Records == nil {
		z.Records = []*ZoneRecord{}
	}
	if z.Filters == nil {
		z.Filters = []*ZoneFilter{}
	}
}

func (z *Zone) AuthKey() string {
	return z.Key()
}

func (z *Zone) SliceOf() interface{} {
	ws := []*Zone{}
	return &ws
}

func (z *Zone) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Zone)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (z *Zone) Validate() {
	z.AddError(ValidName("Invalid Name", z.Name))
	for i := range z.Forwarders {
		s := z.Forwarders[i]

		s = strings.TrimSpace(s)
		if strings.Contains(s, ":") {
			p := strings.Split(s, ":")
			if p[0] == "" {
				z.Errorf("Invalid IP in forwarders list: %s", s)
				continue
			}
			ip := net.ParseIP(p[0])
			if ip == nil {
				z.Errorf("Invalid IP in forwarders list: %s", s)
				continue
			}
			if _, err := strconv.Atoi(p[1]); err != nil {
				z.Errorf("Invalid IP:Port in forwarders list: %s", s)
				continue
			}
		} else {
			ip := net.ParseIP(s)
			if ip == nil {
				z.Errorf("Invalid IP in forwarders list: %s", s)
				continue
			}
			s = s + ":53"
		}
		z.Forwarders[i] = s
	}
	for i, fi := range z.Filters {
		switch fi.Type {
		case 1:
			if _, _, e := net.ParseCIDR(fi.Filter); e != nil {
				z.Errorf("Invalid filter subnet (%s) at index %d: %v", fi.Filter, i, e)
			}
		default:
			z.Errorf("Invalid filter type (%d) at index %d", fi.Type, i)
		}
	}
}

func (z *Zone) CanHaveActions() bool {
	return true
}
