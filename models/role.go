package models

import (
	"fmt"
	"log"
	"strings"

	"github.com/VictorLowther/jsonpatch2"
)

func csm(q string) map[string]struct{} {
	res := map[string]struct{}{}
	for _, p := range strings.Split(q, ",") {
		res[strings.TrimSpace(p)] = struct{}{}
	}
	return res
}

var (
	valScopedActions = csm("update, action")
	basicActions     = csm("list, get, create, delete, actions")

	extraScopes = map[string]string{
		"contents":   "list, get, create, update, delete",
		"files":      "list, get, post, delete",
		"interfaces": "list, get",
		"info":       "get",
		"isos":       "list, get, post, delete",
		"indexes":    "list, get",
	}

	addedActions = map[string]string{
		"users": "token, password",
		"jobs":  "log",
	}

	overriddenActions = map[string]string{
		"prefs":  "list, post",
		"events": "post",
	}

	allScopes = func() map[string]map[string]struct{} {
		res := map[string]map[string]struct{}{}
		for k, v := range extraScopes {
			res[k] = csm(v)
		}
		for _, k := range AllPrefixes() {
			actions := map[string]struct{}{}
			for k2, v2 := range basicActions {
				actions[k2] = v2
			}
			for k2, v2 := range valScopedActions {
				actions[k2] = v2
			}
			if v, ok := addedActions[k]; ok {
				for i := range csm(v) {
					actions[i] = struct{}{}
				}
			}
			if v, ok := overriddenActions[k]; ok {
				actions = csm(v)
			}
			res[k] = actions
		}
		return res
	}()
)

type actionNode struct {
	instances string
}

// actionNode A contains actionNode b if every instance in b is also in a
func (a actionNode) contains(b actionNode) bool {
	// Star case
	if a.instances == "*" {
		return true
	}
	ai, bi := csm(a.instances), csm(b.instances)
	// Adding tenants will make this more complicated.
	for key := range bi {
		if _, ok := ai[key]; !ok {
			return false
		}
	}
	return true
}

type scopeNode struct {
	actions map[string]actionNode
}

func (a scopeNode) contains(b scopeNode) bool {
	for key, ba := range b.actions {
		aa, ok := a.actions[key]
		if !ok {
			parts := strings.SplitN(key, ":", 2)
			if len(parts) == 2 {
				switch parts[0] {
				case "action":
					aa, ok = a.actions[parts[0]]
				case "update":
					ptr, err := jsonpatch2.NewPointer(parts[1])
					for err == nil {
						_, ptr2 := ptr.Chop()
						if len(ptr) == len(ptr2) {
							aa, ok = a.actions[parts[0]]
							break
						}
						frag := parts[0] + ":" + ptr2.String()
						aa, ok = a.actions[frag]
						if ok {
							break
						}
						ptr = ptr2
					}
				}
			}
		}
		if !(ok && aa.contains(ba)) {
			return false
		}
	}
	return true
}

// Claim is an individial specifier for something we are allowed access to.
// User is an API user of DigitalRebar Provision
// swagger:model
type Claim struct {
	Scope    string `json:"scope"`
	Action   string `json:"action"`
	Specific string `json:"specific"`
	scopes   map[string]scopeNode
}

func (c *Claim) compile(e ErrorAdder) {
	c.scopes = map[string]scopeNode{}
	if c.Scope == "*" {
		for k := range allScopes {
			c.scopes[k] = scopeNode{actions: map[string]actionNode{}}
		}
	} else {
		for k := range csm(c.Scope) {
			if _, ok := allScopes[k]; ok {
				c.scopes[k] = scopeNode{actions: map[string]actionNode{}}
			} else if e != nil {
				e.Errorf("No such scope '%s'", k)
			}
		}
	}
	for k := range c.scopes {
		if c.Action == "*" {
			for k2 := range allScopes[k] {
				c.scopes[k].actions[k2] = actionNode{instances: c.Specific}
			}
		} else {
			for k2 := range csm(c.Action) {
				parts := strings.SplitN(k2, ":", 2)
				if _, ok := allScopes[k][parts[0]]; ok {
					c.scopes[k].actions[k2] = actionNode{instances: c.Specific}
				} else if e != nil {
					e.Errorf("No such action '%s'", k2)
				}
			}
		}
	}
}

func (a *Claim) Contains(b *Claim) bool {
	if a.scopes == nil {
		a.compile(nil)
	}
	if b.scopes == nil {
		b.compile(nil)
	}
	for key, bs := range b.scopes {
		as, ok := a.scopes[key]
		if !ok {
			return false
		}
		if !as.contains(bs) {
			return false
		}
	}
	return true
}

// Match tests to see if this claim allows access for the specified
// scope, action, and specific item.
func (c *Claim) Match(scope, action, specific string) bool {
	if c.scopes == nil {
		c.compile(nil)
	}
	c2 := &Claim{Scope: scope, Action: action, Specific: specific}
	c2.compile(nil)
	return c.Contains(c2)
}

func (c *Claim) Validate(e ErrorAdder) {
	c.compile(e)
}

func (c *Claim) String() string {
	return fmt.Sprintf("%s %s %s", c.Scope, c.Action, c.Specific)
}

func MakeClaims(things ...string) []*Claim {
	if len(things)%3 != 0 {
		log.Printf("Bad claim %v", things)
		panic("Strings passed to claims must be a multiple of 3")
	}
	res := make([]*Claim, 0, len(things)/3)
	for len(things) > 0 {
		res = append(res, &Claim{Scope: things[0], Action: things[1], Specific: things[2]})
		things = things[3:]
	}
	return res
}

// User is an API user of DigitalRebar Provision
// swagger:model
type Role struct {
	Validation
	Access
	Meta
	// Name is the name of the user
	//
	// required: true
	Name string
	// Description of user
	Description string
	// Claims that the role support.
	Claims []*Claim
}

func (r *Role) Fill() {
	r.Validation.fill()
	if r.Meta == nil {
		r.Meta = Meta{}
	}
	if r.Claims == nil {
		r.Claims = []*Claim{}
	}
}

// Role a contains role b if a can be used to satisfy all requests b can satisfy
func (a *Role) Contains(b *Role) bool {
	finalRes := true
	res := false
	for _, bClaim := range b.Claims {
		for _, aClaim := range a.Claims {
			res = aClaim.Contains(bClaim)
			if res {
				break
			}
		}
		finalRes = res
		if !finalRes {
			break
		}
	}
	return finalRes
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

func MakeRole(name string, claims ...string) *Role {
	return &Role{
		Name:   name,
		Claims: MakeClaims(claims...),
	}
}
