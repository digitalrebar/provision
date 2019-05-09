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
	// Index of the interface
	//
	Index int
	// A List of Addresses on the interface (CIDR)
	//
	// required: true
	Addresses []string
	// The interface to use for this interface when
	// advertising or claiming access (CIDR)
	//
	ActiveAddress string
	// Possible gateway for this interface
	Gateway string
	// Possible DNS for this interface
	DnsServers []string
	// Possible DNS for domain for this interface
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
