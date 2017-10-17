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
	// Token secret - this is used when generating user token's to
	// allow for revocation by the grantor or the grantee.  Changing this
	// will invalidate all existing tokens that have this user as a user
	// or a grantor.
	Secret string
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

// swagger:model
type UserPassword struct {
	Password string
}

type Users []*User

func (s Users) Elem() Model {
	return &User{}
}

func (s Users) Items() []Model {
	res := make([]Model, len(s))
	for i, m := range s {
		res[i] = m
	}
	return res
}

func (s Users) Fill(m []Model) {
	q := make([]*User, len(m))
	for i, obj := range m {
		q[i] = obj.(*User)
	}
	s = q[:]
}
