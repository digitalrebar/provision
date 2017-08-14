package models

import "github.com/xeipuuv/gojsonschema"

// Param represents metadata about a Parameter or a Preference.
// Specifically, it contains a description of what the information
// is for, detailed documentation about the param, and a JSON schema that
// the param must match to be considered valid.
// swagger:model
type Param struct {
	Validation
	// Name is the name of the param.  Params must be uniquely named.
	//
	// required: true
	Name string
	// Description is a one-line description of the parameter.
	Description string
	// Documentation details what the parameter does, what values it can
	// take, what it is used for, etc.
	Documentation string
	// Schema must be a valid JSONSchema as of draft v4.
	//
	// required: true
	Schema interface{}
}

func (p *Param) Prefix() string {
	return "params"
}

func (p *Param) Key() string {
	return p.Name
}

func (p *Param) ValidateSchema() error {
	_, err := gojsonschema.NewSchema(gojsonschema.NewGoLoader(p.Schema))
	return err
}
