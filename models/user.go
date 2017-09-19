package models

import (
	sc "github.com/elithrar/simple-scrypt"
)

// User is an API user of DigitalRebar Provision
// swagger:model
type User struct {
	Validation
	Access
	MetaData
	// Name is the name of the user
	//
	// required: true
	Name string
	// PasswordHash is the scrypt-hashed version of the user's Password.
	//
	PasswordHash []byte `json:",omitempty"`
}

func (u *User) Prefix() string {
	return "users"
}

func (u *User) Key() string {
	return u.Name
}

func (u *User) CheckPassword(pass string) bool {
	if err := sc.CompareHashAndPassword(u.PasswordHash, []byte(pass)); err == nil {
		return true
	}
	return false
}

func (u *User) Sanitize() Model {
	res, _ := Clone(u)
	res.(*User).PasswordHash = []byte{}
	return res
}

func (u *User) AuthKey() string {
	return u.Key()
}

type Stat struct {
	// required: true
	Name string `json:"name"`
	// required: true
	Count int `json:"count"`
}

// swagger:model
type Info struct {
	// required: true
	Arch string `json:"arch"`
	// required: true
	Os string `json:"os"`
	// required: true
	Version string `json:"version"`
	// required: true
	Id string `json:"id"`
	// required: true
	ApiPort int `json:"api_port"`
	// required: true
	FilePort int `json:"file_port"`
	// required: true
	TftpEnabled bool `json:"tftp_enabled"`
	// required: true
	DhcpEnabled bool `json:"dhcp_enabled"`
	// required: true
	ProvisionerEnabled bool `json:"prov_enabled"`
	// required: true
	Stats []*Stat `json:"stats"`
}

// swagger:model
type UserToken struct {
	Token string
	Info  Info
}

// swagger:model
type UserPassword struct {
	Password string
}
