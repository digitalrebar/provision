package models

// Template represents a template that will be associated with a boot
// environment.
//
// swagger:model
type Template struct {
	Validation
	// ID is a unique identifier for this template.  It cannot change once it is set.
	//
	// required: true
	ID string
	// A description of this template
	Description string
	// Contents is the raw template.  It must be a valid template
	// according to text/template.
	//
	// required: true
	Contents string
}

func (t *Template) Prefix() string {
	return "templates"
}

func (t *Template) Key() string {
	return t.ID
}
