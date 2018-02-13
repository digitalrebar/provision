package models

// Pref tracks a global DigitalRebar Provision preference -- things like the
// bootenv to use for unknown systems trying to PXE boot to us, the
// default bootenv for known systems, etc.
//
type Pref struct {
	Meta
	Name string
	Val  string
}

func (p *Pref) Prefix() string {
	return "preferences"
}

func (p *Pref) Key() string {
	return p.Name
}

func (p *Pref) KeyName() string {
	return "Name"
}

func (p *Pref) Fill() {
	if p.Meta == nil {
		p.Meta = Meta{}
	}
}

func (p *Pref) AuthKey() string {
	return p.Key()
}

func (b *Pref) SliceOf() interface{} {
	s := []*Pref{}
	return &s
}

func (b *Pref) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Pref)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}
