package models

// Access holds if the object is read-only or not
//
// swagger: model
type Access struct {
	// ReadOnly tracks if the store for this object is read-only
	// read only: true
	ReadOnly bool
}

//
// model object may define a Validate method that can
// be used to return errors about if the model is valid
// in the current datatracker.
//
type Accessor interface {
	IsReadOnly() bool
}

func (a *Access) IsReadOnly() bool {
	return a.ReadOnly
}
