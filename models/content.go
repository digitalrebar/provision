package models

type ContentMetaData struct {
	MetaData

	// required: true
	Name        string
	Source      string
	Description string
	Version     string

	// Informational Fields
	Writable     bool
	Type         string
	Overwritable bool
}

//
// Isos???
// Files??
//
// swagger:model
type Content struct {
	// required: true
	Meta ContentMetaData `json:"meta"`

	/*
		These are the sections:
		tasks        map[string]*models.Task
		bootenvs     map[string]*models.BootEnv
		stages       map[string]*models.Stage
		templates    map[string]*models.Template
		profiles     map[string]*models.Profile
		params       map[string]*models.Param
		reservations map[string]*models.Reservation
		subnets      map[string]*models.Subnet
		users        map[string]*models.User
		preferences  map[string]*models.Pref
		plugins      map[string]*models.Plugin
		machines     map[string]*models.Machine
		leases       map[string]*models.Lease
	*/
	Sections Sections `json:"sections"`
}

type Sections map[string]Section
type Section map[string]interface{}

// swagger:model
type ContentSummary struct {
	Meta     ContentMetaData `json:"meta"`
	Counts   map[string]int
	Warnings []string
}

func (c *Content) Prefix() string {
	return "contents"
}

func (c *Content) Key() string {
	return c.Meta.Name
}

func (c *Content) AuthKey() string {
	return c.Key()
}

type Contents []*Content

func (s Contents) Elem() Model {
	return &Content{}
}

func (s Contents) Items() []Model {
	res := make([]Model, len(s))
	for i, m := range s {
		res[i] = m
	}
	return res
}
func (s Contents) Fill(m []Model) {
	q := make([]*Content, len(m))
	for i, obj := range m {
		q[i] = obj.(*Content)
	}
	s = q[:]
}
