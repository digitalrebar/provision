package models

// Owned holds the info about which DRP Endpoint owns this object
//
// swagger: model
type Owned struct {
	// Endpoint tracks the owner of the object amoung DRP endpoints
	// read only: true
	Endpoint string
}

//
// model object may define a GetEndpoint() method that can
// be used to return the owner for the object
//
type Owner interface {
	GetEndpoint() string
}

// GetEndpoint returns the name of the owning DRP Endpoint
func (o *Owned) GetEndpoint() string {
	return o.Endpoint
}
