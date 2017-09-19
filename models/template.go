package models

// Template represents a template that will be associated with a boot
// environment.
//
// swagger:model
type Template struct {
	Validation
	Access
	MetaData
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

func (t *Template) AuthKey() string {
	return t.Key()
}

type Templates []*Template

func (s Templates) Elem() Model {
	return &Template{}
}

func (s Templates) Items() []Model {
	res := make([]Model, len(s))
	for i, m := range s {
		res[i] = m
	}
	return res
}

func (s Templates) Fill(m []Model) {
	q := make([]*Template, len(m))
	for i, obj := range m {
		q[i] = obj.(*Template)
	}
	s = q[:]
}
