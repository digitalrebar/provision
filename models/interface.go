package models

// swagger:model
type Interface struct {
	Access
	MetaData
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
}

func (i *Interface) Prefix() string { return "interfaces" }
func (i *Interface) Key() string    { return i.Name }

type Interfaces []*Interface

func (s Interfaces) Elem() Model {
	return &Interface{}
}

func (s Interfaces) Items() []Model {
	res := make([]Model, len(s))
	for i, m := range s {
		res[i] = m
	}
	return res
}

func (s Interfaces) Fill(m []Model) {
	q := make([]*Interface, len(m))
	for i, obj := range m {
		q[i] = obj.(*Interface)
	}
	s = q[:]
}
