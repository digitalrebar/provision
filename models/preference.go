package models

// Pref tracks a global DigitalRebar Provision preference -- things like the
// bootenv to use for unknown systems trying to PXE boot to us, the
// default bootenv for known systems, etc.
//
type Pref struct {
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
