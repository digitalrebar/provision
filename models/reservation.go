package models

import "net"

// Reservation tracks persistent DHCP IP address reservations.
//
// swagger:model
type Reservation struct {
	Validation
	Access
	MetaData
	// Addr is the IP address permanently assigned to the strategy/token combination.
	//
	// required: true
	// swagger:strfmt ipv4
	Addr net.IP
	// Token is the unique identifier that the strategy for this Reservation should use.
	//
	// required: true
	Token string
	// NextServer is the address the server should contact next.
	//
	// required: false
	// swagger:strfmt ipv4
	NextServer net.IP
	// Options is the list of DHCP options that apply to this Reservation
	Options []DhcpOption
	// Strategy is the leasing strategy that will be used determine what to use from
	// the DHCP packet to handle lease management.
	//
	// required: true
	Strategy string
}

func (r *Reservation) Prefix() string {
	return "reservations"
}

func (r *Reservation) Key() string {
	return Hexaddr(r.Addr)
}

func (r *Reservation) AuthKey() string {
	return r.Key()
}

type Reservations []*Reservation

func (s Reservations) Elem() Model {
	return &Reservation{}
}

func (s Reservations) Items() []Model {
	res := make([]Model, len(s))
	for i, m := range s {
		res[i] = m
	}
	return res
}

func (s Reservations) Fill(m []Model) {
	q := make([]*Reservation, len(m))
	for i, obj := range m {
		q[i] = obj.(*Reservation)
	}
	s = q[:]
}
