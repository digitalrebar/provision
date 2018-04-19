package models

import "strings"

func csp(v string) []string {
	res := strings.Split(v, ",")
	for i := range res {
		res[i] = strings.TrimSpace(res[i])
	}
	return res
}

var validActions = map[string]struct{}{
	"actions":  struct{}{},
	"create":   struct{}{},
	"delete":   struct{}{},
	"get":      struct{}{},
	"list":     struct{}{},
	"log":      struct{}{},
	"password": struct{}{},
	"patch":    struct{}{},
	"post":     struct{}{},
	"token":    struct{}{},
	"update":   struct{}{},
}

// Claim is an individial specifier for something we are allowed access to.
type Claim struct {
	Scope    string `json:"scope"`
	Action   string `json:"action"`
	Specific string `json:"specific"`
}

// Match tests to see if this claim allows access for the specified
// scope, action, and specific item.
//
// If the Claim has `*` for any field, it matches all possible values
// for that field.
func (c *Claim) Match(scope, action, specific string) bool {
	scopeMatch, actionMatch, specificMatch := c.Scope == "*", c.Action == "*", c.Specific == "*"
	if !scopeMatch {
		for _, sc := range csp(c.Scope) {
			scopeMatch = sc == scope
			if scopeMatch {
				break
			}
		}
	}
	if !actionMatch {
		for _, ac := range csp(c.Action) {
			actionMatch = ac == action
			if actionMatch {
				break
			}
		}
	}
	if !specificMatch {
		for _, sc := range csp(c.Specific) {
			specificMatch = sc == specific
			if specificMatch {
				break
			}
		}
	}
	return scopeMatch && actionMatch && specificMatch
}

func (c *Claim) Validate(r *Role) {
	if c.Scope != "*" {
		for _, sc := range csp(c.Scope) {
			r.AddError(ValidName("Invalid Scope", sc))
		}
	}
	if c.Action != "*" {
		for _, sc := range csp(c.Action) {
			if _, ok := validActions[sc]; !ok {
				r.Errorf("Invalid Action %s for claim scope %s", sc, c.Scope)
			}
		}
	}
	if c.Specific != "*" {
		for _, sc := range csp(c.Specific) {
			r.AddError(ValidName("Invalid Specific", sc))
		}
	}
}

type Role struct {
	Validation
	Access
	Meta
	Name   string
	Claims []Claim
}

func (r *Role) Fill() {
	r.Validation.fill()
	if r.Meta == nil {
		r.Meta = Meta{}
	}
	if r.Claims == nil {
		r.Claims = []Claim{}
	}
}

func (r *Role) Validate() {
	r.AddError(ValidName("Invalid Name", r.Name))
	for _, c := range r.Claims {
		c.Validate(r)
	}
}

func (r *Role) Prefix() string {
	return "roles"
}

func (r *Role) Key() string {
	return r.Name
}

func (r *Role) KeyName() string {
	return "Name"
}

func (r *Role) AuthKey() string {
	return r.Key()
}

func (r *Role) SliceOf() interface{} {
	rs := []*Role{}
	return &rs
}

func (r *Role) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Role)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (r *Role) Match(scope, action, specific string) bool {
	for _, c := range r.Claims {
		if c.Match(scope, action, specific) {
			return true
		}
	}
	return false
}
