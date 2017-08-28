package models

// Access holds if the object is read-only or not
//
// swagger: model
type Access struct {
	// ReadOnly tracks if the store for this object is read-only
	// read only: true
	ReadOnly bool
}
