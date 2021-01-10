package models

// Interface represents a network interface that is present on the
// server running dr-provision.  It is primarily used by the UX to
// help generate Subnets.
//
// swagger:model
type Interface struct {
	Access
	Meta
	// Name of the interface
	//
	// required: true
	Name string
	// Index of the interface.  This is OS specific.
	Index int
	// Addresses contains the IPv4 and IPv6 addresses bound to this interface in no particular order.
	//
	// required: true
	Addresses []string
	// ActiveAddress is our best guess at the address that should be used for "normal" incoming traffic
	// on this interface.
	ActiveAddress string
	// Gateway is our best guess about the IP address that traffic forwarded through this interface should
	// be sent to.
	Gateway string
	// DnsServers is a list of DNS server that hsould be used when resolving addresses via this interface.
	DnsServers []string
	// DnsDomain is the domain that this system appears to be in on this interface.
	DnsDomain string
}

func (i *Interface) Prefix() string  { return "interfaces" }
func (i *Interface) Key() string     { return i.Name }
func (i *Interface) KeyName() string { return "Name" }
func (i *Interface) Fill() {
	if i.Meta == nil {
		i.Meta = Meta{}
	}
	if i.Addresses == nil {
		i.Addresses = []string{}
	}
}

func (b *Interface) SliceOf() interface{} {
	s := []*Interface{}
	return &s
}

func (b *Interface) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Interface)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}
