package models

// Partialed holds if the object is partially filled in.
//
// swagger: model
type Partialed struct {
	// Partial tracks if the object is not complete when returned.
	// read only: true
	Partial bool
}

// Partialer is an interface that objects that are partially returned.
type Partialer interface {
	IsPartial() bool
	SetPartial()
}

// IsPartial returns whether the object is partially returned.
// This will be set if the object has been slimmed or partially returned.
func (p *Partialed) IsPartial() bool {
	return p.Partial
}

func (p *Partialed) SetPartial() {
	p.Partial = true
}
