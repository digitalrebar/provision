package store

// When loading an object, the store should
// set the current readonly state on the object.
type ReadOnlySetter interface {
	SetReadOnly(bool)
}

// When loading an object, the store should
// set the current owner name.
type BundleSetter interface {
	SetBundle(string)
}
