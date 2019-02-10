package models

// Template represents a template that will be associated with a boot
// environment.
//
// swagger:model
type Template struct {
	Validation
	Access
	Meta
	Owned
	Bundled
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

func (t *Template) GetMeta() Meta {
	return t.Meta
}

func (t *Template) SetMeta(d Meta) {
	t.Meta = d
}

func (t *Template) Validate() {
	t.AddError(ValidName("Invalid ID", t.ID))
}

func (t *Template) Prefix() string {
	return "templates"
}

func (t *Template) Key() string {
	return t.ID
}

func (t *Template) KeyName() string {
	return "ID"
}

func (t *Template) Fill() {
	t.Validation.fill()
	if t.Meta == nil {
		t.Meta = Meta{}
	}
}

func (t *Template) AuthKey() string {
	return t.Key()
}

func (b *Template) SliceOf() interface{} {
	s := []*Template{}
	return &s
}

func (b *Template) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Template)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (b *Template) CanHaveActions() bool {
	return true
}
