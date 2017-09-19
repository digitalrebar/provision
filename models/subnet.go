package models

import "net"

// Subnet represents a DHCP Subnet
//
// swagger:model
type Subnet struct {
	Validation
	Access
	MetaData
	// Name is the name of the subnet.
	// Subnet names must be unique
	//
	// required: true
	Name string
	// Enabled indicates if the subnet should hand out leases or continue operating
	// leases if already running.
	//
	// required: true
	Enabled bool
	// Subnet is the network address in CIDR form that all leases
	// acquired in its range will use for options, lease times, and NextServer settings
	// by default
	//
	// required: true
	// pattern: ^([0-9]+\.){3}[0-9]+/[0-9]+$
	Subnet string
	// NextServer is the address of the next server
	//
	// required: true
	// swagger:strfmt ipv4
	NextServer net.IP
	// ActiveStart is the first non-reserved IP address we will hand
	// non-reserved leases from.
	//
	// required: true
	// swagger:strfmt ipv4
	ActiveStart net.IP
	// ActiveEnd is the last non-reserved IP address we will hand
	// non-reserved leases from.
	//
	// required: true
	// swagger:strfmt ipv4
	ActiveEnd net.IP
	// ActiveLeaseTime is the default lease duration in seconds
	// we will hand out to leases that do not have a reservation.
	//
	// required: true
	ActiveLeaseTime int32
	// ReservedLeasTime is the default lease time we will hand out
	// to leases created from a reservation in our subnet.
	//
	// required: true
	ReservedLeaseTime int32
	// OnlyReservations indicates that we will only allow leases for which
	// there is a preexisting reservation.
	//
	// required: true
	OnlyReservations bool
	Options          []*DhcpOption
	// Strategy is the leasing strategy that will be used determine what to use from
	// the DHCP packet to handle lease management.
	//
	// required: true
	Strategy string
	// Pickers is list of methods that will allocate IP addresses.
	// Each string must refer to a valid address picking strategy.  The current ones are:
	//
	// "none", which will refuse to hand out an address and refuse
	// to try any remaining strategies.
	//
	// "hint", which will try to reuse the address that the DHCP
	// packet is requesting, if it has one.  If the request does
	// not have a requested address, "hint" will fall through to
	// the next strategy. Otherwise, it will refuse to try ant
	// reamining strategies whether or not it can satisfy the
	// request.  This should force the client to fall back to
	// DHCPDISCOVER with no requsted IP address. "hint" will reuse
	// expired leases and unexpired leases that match on the
	// requested address, strategy, and token.
	//
	// "nextFree", which will try to create a Lease with the next
	// free address in the subnet active range.  It will fall
	// through to the next strategy if it cannot find a free IP.
	// "nextFree" only considers addresses that do not have a
	// lease, whether or not the lease is expired.
	//
	// "mostExpired" will try to recycle the most expired lease in the subnet's active range.
	//
	// All of the address allocation strategies do not consider
	// any addresses that are reserved, as lease creation will be
	// handled by the reservation instead.
	//
	// We will consider adding more address allocation strategies in the future.
	//
	// required: true
	Pickers []string
}

func (s *Subnet) Prefix() string {
	return "subnets"
}

func (s *Subnet) Key() string {
	return s.Name
}

func (s *Subnet) AuthKey() string {
	return s.Key()
}

type Subnets []*Subnet

func (s Subnets) Elem() Model {
	return &Subnet{}
}

func (s Subnets) Items() []Model {
	res := make([]Model, len(s))
	for i, m := range s {
		res[i] = m
	}
	return res
}

func (s Subnets) Fill(m []Model) {
	q := make([]*Subnet, len(m))
	for i, obj := range m {
		q[i] = obj.(*Subnet)
	}
	s = q[:]
}
