package models

import (
	"fmt"
	"net"
	"time"
)

var hexDigit = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F'}

func Hexaddr(addr net.IP) string {
	b := addr.To4()
	s := make([]byte, len(b)*2)
	for i, tn := range b {
		s[i*2], s[i*2+1] = hexDigit[tn>>4], hexDigit[tn&0xf]
	}
	return string(s)
}

// Lease tracks DHCP leases.
// swagger:model
type Lease struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	// Addr is the IP address that the lease handed out.
	//
	// required: true
	// swagger:strfmt ipv4
	Addr net.IP
	// NextServer is the IP address that we should have the machine talk to
	// next.  In most cases, this will be our address.
	//
	// required: false
	// swagger:strfmt ipv4
	NextServer net.IP
	// Via is the IP address used to select which subnet the lease belongs to.
	// It is either an address present on a local interface that dr-provision is
	// listening on, or the GIADDR field of the DHCP request.
	//
	// required: false
	// swagger:strfmt ipv4
	Via net.IP
	// Token is the unique token for this lease based on the
	// Strategy this lease used.
	//
	// required: true
	Token string
	// Duration is the time in seconds for which a lease can be valid.
	// ExpireTime is calculated from Duration.
	Duration int32
	// ExpireTime is the time at which the lease expires and is no
	// longer valid The DHCP renewal time will be half this, and the
	// DHCP rebind time will be three quarters of this.
	//
	// required: true
	// swagger:strfmt date-time
	ExpireTime time.Time
	// Strategy is the leasing strategy that will be used determine what to use from
	// the DHCP packet to handle lease management.
	//
	// required: true
	Strategy string
	// State is the current state of the lease.  This field is for informational
	// purposes only.
	//
	// read only: true
	// required: true
	State string
	// Options are the DHCP options that the Lease is running with.
	Options []DhcpOption
	// ProvidedOptions are the DHCP options the last Discover or Offer packet
	// for this lease provided to us.
	ProvidedOptions []DhcpOption
	// SkipBoot indicates that the DHCP system is allowed to offer
	// boot options for whatever boot protocol the machine wants to
	// use.
	//
	// read only: true
	SkipBoot bool
}

func (l *Lease) String() string {
	return fmt.Sprintf("%s %s:%s:%s %d", l.Addr, l.Strategy, l.Token, l.State, l.ExpireTime.Unix())
}

func (l *Lease) GetMeta() Meta {
	return l.Meta
}

func (l *Lease) SetMeta(d Meta) {
	l.Meta = d
}

func (l *Lease) Prefix() string {
	return "leases"
}

func (l *Lease) Key() string {
	return Hexaddr(l.Addr)
}

func (l *Lease) KeyName() string {
	return "Addr"
}

func (l *Lease) Fill() {
	if l.Meta == nil {
		l.Meta = Meta{}
	}
	if l.NextServer == nil {
		l.NextServer = net.IP{}
	}
	if l.Via == nil {
		l.Via = net.IP{}
	}
	if l.Options == nil {
		l.Options = []DhcpOption{}
	}
	if l.ProvidedOptions == nil {
		l.ProvidedOptions = []DhcpOption{}
	}
	l.Validation.fill()
}

func (l *Lease) AuthKey() string {
	return l.Key()
}

func (b *Lease) SliceOf() interface{} {
	s := []*Lease{}
	return &s
}

func (b *Lease) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Lease)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (b *Lease) CanHaveActions() bool {
	return true
}

func (l *Lease) Expired() bool {
	return l.ExpireTime.Before(time.Now())
}

func (l *Lease) Fake() bool {
	return l.State == "FAKE"
}

func (l *Lease) Expire() {
	l.ExpireTime = time.Now()
	l.State = "EXPIRED"
}

func (l *Lease) Invalidate() {
	l.ExpireTime = time.Now().Add(10 * time.Minute)
	l.Token = ""
	l.Strategy = ""
	l.State = "INVALID"
}
