package models

import (
	"math/big"
	"net"
)

// Subnet represents a DHCP Subnet
//
// swagger:model
type Subnet struct {
	Validation
	Access
	Meta
	Owned
	// Name is the name of the subnet.
	// Subnet names must be unique
	//
	// required: true
	Name string
	// A description of this Subnet.  This should tell what it is for,
	// any special considerations that should be taken into account when
	// using it, etc.
	Description string
	// Documentation of this subnet.  This should tell what
	// the subnet is for, any special considerations that
	// should be taken into account when using it, etc. in rich structured text (rst).
	Documentation string
	// Enabled indicates if the subnet should hand out leases or continue operating
	// leases if already running.
	//
	// required: true
	Enabled bool
	// Proxy indicates if the subnet should act as a proxy DHCP server.
	// If true, the subnet will not manage ip addresses but will send
	// offers to requests.  It is an error for Proxy and Unmanaged to be
	// true.
	//
	// required: true
	Proxy bool
	// Unmanaged indicates that dr-provision will never send
	// boot-related options to machines that get leases from this
	// subnet.  If false, dr-provision will send whatever boot-related
	// options it would normally send.  It is an error for Unmanaged and
	// Proxy to both be true.
	//
	// required: true
	Unmanaged bool
	// Subnet is the network address in CIDR form that all leases
	// acquired in its range will use for options, lease times, and NextServer settings
	// by default
	//
	// required: true
	// pattern: ^([0-9]+\.){3}[0-9]+/[0-9]+$
	Subnet string
	// NextServer is the address of the next server in the DHCP/TFTP/PXE
	// chain.  You should only set this if you want to transfer control
	// to a different DHCP or TFTP server.
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
	Options          []DhcpOption
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
	// the next strategy. Otherwise, it will refuse to try any
	// remaining strategies whether or not it can satisfy the
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

func (s *Subnet) GetMeta() Meta {
	return s.Meta
}

func (s *Subnet) SetMeta(d Meta) {
	s.Meta = d
}

func (s *Subnet) GetDocumentation() string {
	return s.Documentation
}

func (s *Subnet) Validate() {
	s.AddError(ValidName("Invalid Name", s.Name))
	_, subnet, err := net.ParseCIDR(s.Subnet)
	if err != nil {
		s.Errorf("Invalid subnet %s: %v", s.Subnet, err)
		return
	} else {
		ValidateIP4(s, subnet.IP)
	}
	if s.Strategy == "" {
		s.Errorf("Strategy must have a value")
	}
	if s.NextServer != nil {
		ValidateMaybeZeroIP4(s, s.NextServer)
	}
	if s.Proxy && s.Unmanaged {
		s.Errorf("Unmanaged and Proxy cannot both be true")
	}
	if !(s.OnlyReservations || s.Proxy) {
		ValidateIP4(s, s.ActiveStart)
		ValidateIP4(s, s.ActiveEnd)
		if !subnet.Contains(s.ActiveStart) {
			s.Errorf("ActiveStart %s not in subnet range %s", s.ActiveStart, subnet)
		}
		if !subnet.Contains(s.ActiveEnd) {
			s.Errorf("ActiveEnd %s not in subnet range %s", s.ActiveEnd, subnet)
		}
		startBytes := big.NewInt(0)
		endBytes := big.NewInt(0)
		startBytes.SetBytes(s.ActiveStart)
		endBytes.SetBytes(s.ActiveEnd)
		if startBytes.Cmp(endBytes) != -1 {
			s.Errorf("ActiveStart %s must be less than ActiveEnd %s", s.ActiveStart, s.ActiveEnd)
		}
		if s.ActiveLeaseTime < 60 {
			s.Errorf("ActiveLeaseTime must be greater than or equal to 60 seconds, not %d", s.ActiveLeaseTime)
		}
	}
	if s.ReservedLeaseTime < 7200 {
		s.Errorf("ReservedLeaseTime must be greater than or equal to 7200 seconds, not %d", s.ReservedLeaseTime)
	}

}

func (s *Subnet) Prefix() string {
	return "subnets"
}

func (s *Subnet) Key() string {
	return s.Name
}

func (s *Subnet) KeyName() string {
	return "Name"
}

func (s *Subnet) Fill() {
	s.Validation.fill()
	if s.Meta == nil {
		s.Meta = Meta{}
	}
	if s.Options == nil {
		s.Options = []DhcpOption{}
	}
	if s.Strategy == "" {
		s.Strategy = "MAC"
	}
	if s.Pickers == nil || len(s.Pickers) == 0 {
		if s.OnlyReservations {
			s.Pickers = []string{"none"}
		} else {
			s.Pickers = []string{"hint", "nextFree", "mostExpired"}
		}
	}
	if s.ActiveLeaseTime == 0 {
		s.ActiveLeaseTime = 60
	}
	if s.ReservedLeaseTime == 0 {
		s.ReservedLeaseTime = 7200
	}
}

func (s *Subnet) AuthKey() string {
	return s.Key()
}

func (b *Subnet) SliceOf() interface{} {
	s := []*Subnet{}
	return &s
}

func (b *Subnet) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Subnet)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (b *Subnet) CanHaveActions() bool {
	return true
}
