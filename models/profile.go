package models

// Profile represents a set of key/values to use in
// template expansion.
//
// There is one special profile named 'global' that acts
// as a global set of parameters for the system.
//
// These can be assigned to a machine's profile list.
// swagger:model
type Profile struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	// The name of the profile.  This must be unique across all
	// profiles.
	//
	// required: true
	Name string
	// A description of this profile.  This can contain any reference
	// information for humans you want associated with the profile.
	Description string
	// Documentation of this profile.  This should tell what
	// the profile is for, any special considerations that
	// should be taken into account when using it, etc. in rich structured text (rst).
	Documentation string
	// Any additional parameters that may be needed to expand templates
	// for BootEnv, as documented by that boot environment's
	// RequiredParams and OptionalParams.
	Params map[string]interface{}
}

func (p *Profile) GetMeta() Meta {
	return p.Meta
}

func (p *Profile) SetMeta(d Meta) {
	p.Meta = d
}

func (p *Profile) Validate() {
	p.AddError(ValidName("Invalid Name", p.Name))
	for k := range p.Params {
		p.AddError(ValidParamName("Invalid Param Name", k))
	}
}

func (p *Profile) Prefix() string {
	return "profiles"
}

func (p *Profile) Key() string {
	return p.Name
}

func (p *Profile) KeyName() string {
	return "Name"
}

func (p *Profile) Fill() {
	p.Validation.fill()
	if p.Meta == nil {
		p.Meta = Meta{}
	}
	if p.Params == nil {
		p.Params = map[string]interface{}{}
	}
}

func (p *Profile) AuthKey() string {
	return p.Key()
}

func (p *Profile) SliceOf() interface{} {
	s := []*Profile{}
	return &s
}

func (p *Profile) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Profile)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

// match Paramer interface
func (p *Profile) GetParams() map[string]interface{} {
	return copyMap(p.Params)
}

func (p *Profile) SetParams(pl map[string]interface{}) {
	p.Params = copyMap(pl)
}

func (p *Profile) SetName(n string) {
	p.Name = n
}

func (p *Profile) CanHaveActions() bool {
	return true
}
