package models

// swagger:model
type Interface struct {
	Access
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
