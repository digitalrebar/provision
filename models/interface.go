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
func (i *Interface) Fill() {
	i.MetaData.fill()
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
