package models

// swagger:model
type Tenant struct {
	Validation
	Access
	Meta
	Name        string
	Description string
	Members     map[string][]string
	Users       []string
}

func (t *Tenant) Fill() {
	t.Validation.fill()
	if t.Meta == nil {
		t.Meta = Meta{}
	}
	if t.Members == nil {
		t.Members = map[string][]string{}
	}
	if t.Users == nil {
		t.Users = []string{}
	}
}

func (t *Tenant) GetMeta() Meta {
	return t.Meta
}

func (t *Tenant) SetMeta(d Meta) {
	t.Meta = d
}

func (t *Tenant) Validate() {
	t.AddError(ValidName("Invalid Name", t.Name))
	for k := range t.Members {
		if _, ok := modelPrefixes[k]; !ok {
			t.Errorf("Invalid ")
		}
	}
}

func (t *Tenant) Prefix() string {
	return "tenants"
}

func (t *Tenant) Key() string {
	return t.Name
}

func (t *Tenant) KeyName() string {
	return "Name"
}

func (t *Tenant) AuthKey() string {
	return t.Key()
}

func (t *Tenant) SliceOf() interface{} {
	ts := []*Tenant{}
	return &ts
}

func (t *Tenant) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Tenant)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}
