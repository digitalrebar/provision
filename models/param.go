package models

import "github.com/xeipuuv/gojsonschema"

// Param represents metadata about a Parameter or a Preference.
// Specifically, it contains a description of what the information
// is for, detailed documentation about the param, and a JSON schema that
// the param must match to be considered valid.
// swagger:model
type Param struct {
	Validation
	Access
	Meta
	// Name is the name of the param.  Params must be uniquely named.
	//
	// required: true
	Name string
	// Description is a one-line description of the parameter.
	Description string
	// Documentation details what the parameter does, what values it can
	// take, what it is used for, etc.
	Documentation string
	// Secure implies that any API interactions with this Param
	// will deal with SecureData values.
	//
	// required: true
	Secure bool
	// Schema must be a valid JSONSchema as of draft v4.
	//1
	// required: true
	Schema interface{}
}

func (p *Param) DefaultValue() (interface{}, bool) {
	if km, ok := p.Schema.(map[string]interface{}); ok {
		v, vok := km["default"]
		return v, vok
	}
	return nil, false
}

func (p *Param) TypeValue() (interface{}, bool) {
	if km, ok := p.Schema.(map[string]interface{}); ok {
		v, vok := km["type"]
		return v, vok
	}
	return nil, false
}

func (p *Param) Validate() {
	p.AddError(ValidParamName("Invalid Name", p.Name))
	if p.Schema != nil {
		validator, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(p.Schema))
		if err != nil {
			p.AddError(err)
			return
		}
		dv, ok := p.DefaultValue()
		if !ok {
			return
		}
		res, err := validator.Validate(gojsonschema.NewGoLoader(dv))
		if err != nil {
			p.Errorf("Error validating default value: %v", err)
		} else if !res.Valid() {
			for _, e := range res.Errors() {
				p.Errorf("Error in default value: %v", e.String())
			}
		}
	}
}

func (p *Param) SetName(s string) {
	p.Name = s
}

func (p *Param) Prefix() string {
	return "params"
}

func (p *Param) Key() string {
	return p.Name
}

func (p *Param) KeyName() string {
	return "Name"
}

func (p *Param) Fill() {
	if p.Meta == nil {
		p.Meta = Meta{}
	}
	p.Validation.fill()
}

func (p *Param) AuthKey() string {
	return p.Key()
}

func (b *Param) SliceOf() interface{} {
	s := []*Param{}
	return &s
}

func (b *Param) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Param)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}
