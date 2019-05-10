package models

// UserToken is an auth token for a specific User.
// The Token section can be used for bearer authentication.
//
// swagger:model
type UserToken struct {
	Token string
	Info  Info
}
