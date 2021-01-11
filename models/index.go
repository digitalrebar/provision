package models

// Index holds details on the index
// swagger:model
type Index struct {
	// Type gives you a rough idea of how the string used to query
	// this index should be formatted.
	Type string
	// Unique tells you whether there can be multiple entries in the
	// index for the same key that refer to different items.
	Unique bool
	// Unordered tells you whether this index cannot be sorted.
	Unordered bool
	// Regex indicates whether you can use the Re filter with this index
	Regex bool
}
