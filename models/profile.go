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
	Partialed
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
	// Additional Profiles that should be considered for parameters
	Profiles []string
}

// GetMeta returns the meta data for this profile
func (p *Profile) GetMeta() Meta {
	return p.Meta
}

// SetMeta sets the meta data for this profile
func (p *Profile) SetMeta(d Meta) {
	p.Meta = d
}

// Validate makes sure that the object is valid (outside of references)
func (p *Profile) Validate() {
	p.AddError(ValidName("Invalid Name", p.Name))
	for k := range p.Params {
		p.AddError(ValidParamName("Invalid Param Name", k))
	}
	for _, v := range p.Profiles {
		p.AddError(ValidName("Invalid Profile Name", v))
	}
}

// Prefix returns the object type
func (p *Profile) Prefix() string {
	return "profiles"
}

// Key returns the primary index for this object
func (p *Profile) Key() string {
	return p.Name
}

// KeyName returns the field of the object that is used as the primary key
func (p *Profile) KeyName() string {
	return "Name"
}

// GetDocumentation returns the object's documentation
func (p *Profile) GetDocumentation() string {
	return p.Documentation
}

// GetDescription returns the object's description
func (p *Profile) GetDescription() string {
	return p.Description
}

// Fill initializes the object
func (p *Profile) Fill() {
	p.Validation.fill(p)
	if p.Meta == nil {
		p.Meta = Meta{}
	}
	if p.Params == nil {
		p.Params = map[string]interface{}{}
	}
	if p.Profiles == nil {
		p.Profiles = []string{}
	}
}

// AuthKey returns the value that should be validated against claims
func (p *Profile) AuthKey() string {
	return p.Key()
}

// SliceOf returns an empty slice of this type of objects
func (p *Profile) SliceOf() interface{} {
	s := []*Profile{}
	return &s
}

// ToModels converts a slice of these specific objects to a slice of Model interfaces
func (p *Profile) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Profile)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

// GetParams returns the current parameters for this profile
// matches Paramer interface
func (p *Profile) GetParams() map[string]interface{} {
	return copyMap(p.Params)
}

// SetParams sets the current parameters for this profile
// matches Paramer interface
func (p *Profile) SetParams(pl map[string]interface{}) {
	p.Params = copyMap(pl)
}

// SetName changes the name of the profile
func (p *Profile) SetName(n string) {
	p.Name = n
}

// CanHaveActions indicates if the object is allowed to have actions
func (p *Profile) CanHaveActions() bool {
	return true
}

// match Profiler interface

// GetProfiles returns the profiles on this profile
func (p *Profile) GetProfiles() []string {
	return p.Profiles
}

// SetProfiles sets the profiles on this profile
func (p *Profile) SetProfiles(np []string) {
	p.Profiles = np
}
