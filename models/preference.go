package models

// Pref tracks a global DigitalRebar Provision preference -- things like the
// bootenv to use for unknown systems trying to PXE boot to us, the
// default bootenv for known systems, etc.
//
type Pref struct {
	MetaData
	Name string
	Val  string
}

func (p *Pref) Prefix() string {
	return "preferences"
}

func (p *Pref) Key() string {
	return p.Name
}

func (p *Pref) AuthKey() string {
	return p.Key()
}

type Prefs []*Pref

func (s Prefs) Elem() Model {
	return &Pref{}
}

func (s Prefs) Items() []Model {
	res := make([]Model, len(s))
	for i, m := range s {
		res[i] = m
	}
	return res
}

func (s Prefs) Fill(m []Model) {
	q := make([]*Pref, len(m))
	for i, obj := range m {
		q[i] = obj.(*Pref)
	}
	s = q[:]
}
