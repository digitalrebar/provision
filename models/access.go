package models

// Access holds if the object is read-only or not
//
// swagger: model
type Access struct {
	// ReadOnly tracks if the store for this object is read-only
	// read only: true
	ReadOnly bool
}

// Accessor is an interface that objects that can be ReadOnly should
// satisfy.  model object may define a Validate method that can be
// used to return errors about if the model is valid in the current
// datatracker.
//
type Accessor interface {
	IsReadOnly() bool
}

// IsReadOnly returns whether the object is read-only.
// This will be set if the object comes from any content layer other
// than the working one (provided by a plugin or a content bundle, etc.)
func (a *Access) IsReadOnly() bool {
	return a.ReadOnly
}
