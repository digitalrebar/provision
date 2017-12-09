package models

import "net"

// Reservation tracks persistent DHCP IP address reservations.
//
// swagger:model
type Reservation struct {
	Validation
	Access
	Meta
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

func (r *Reservation) Fill() {
	r.Validation.fill()
	if r.Meta == nil {
		r.Meta = Meta{}
	}
	if r.Options == nil {
		r.Options = []DhcpOption{}
	}
}

func (r *Reservation) AuthKey() string {
	return r.Key()
}

func (b *Reservation) SliceOf() interface{} {
	s := []*Reservation{}
	return &s
}

func (b *Reservation) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Reservation)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}
