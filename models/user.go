package models

import (
	sc "github.com/elithrar/simple-scrypt"
)

// ExtAuthData is used by external auth systems to pass
// back information about the authenticated user.
type ExtAuthData struct {
	User    *User
	Tenants []string
}

// User is an API user of DigitalRebar Provision
// swagger:model
type User struct {
	Validation
	Access
	Meta
	Owned
	// Name is the name of the user
	//
	// required: true
	Name string
	// Description of user
	Description string
	// PasswordHash is the scrypt-hashed version of the user's Password.
	//
	PasswordHash []byte `json:",omitempty"`
	// Token secret - this is used when generating user token's to
	// allow for revocation by the grantor or the grantee.  Changing this
	// will invalidate all existing tokens that have this user as a user
	// or a grantor.
	Secret string
	// Roles is a list of Roles this User has.
	//
	Roles []string
}

func (u *User) GetMeta() Meta {
	return u.Meta
}

func (u *User) SetMeta(d Meta) {
	u.Meta = d
}

func (u *User) Validate() {
	u.AddError(ValidName("Invalid Name", u.Name))
}

func (u *User) Prefix() string {
	return "users"
}

func (u *User) Key() string {
	return u.Name
}

func (u *User) KeyName() string {
	return "Name"
}

func (u *User) Fill() {
	u.Validation.fill()
	if u.Meta == nil {
		u.Meta = Meta{}
	}
	if u.Roles == nil {
		u.Roles = []string{}
	}
}

func (u *User) CheckPassword(pass string) bool {
	if err := sc.CompareHashAndPassword(u.PasswordHash, []byte(pass)); err == nil {
		return true
	}
	return false
}

func (u *User) ChangePassword(newPass string) error {
	ph, err := sc.GenerateFromPassword([]byte(newPass), sc.DefaultParams)
	if err != nil {
		return err
	}
	u.PasswordHash = ph
	// When a user changes their password, invalidate any previous cached auth tokens.
	u.Secret = RandString(16)
	return nil
}

func (u *User) Sanitize() Model {
	res := Clone(u)
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

func (b *User) SliceOf() interface{} {
	s := []*User{}
	return &s
}

func (b *User) ToModels(obj interface{}) []Model {
	items := obj.(*[]*User)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (b *User) SetName(n string) {
	b.Name = n
}

func (b *User) CanHaveActions() bool {
	return true
}
