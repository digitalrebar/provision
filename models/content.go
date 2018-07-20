package models

import "github.com/digitalrebar/store"

// All fields must be strings
type ContentMetaData struct {
	// required: true
	Name        string
	Source      string
	Description string
	Version     string

	// Optional fields
	Documentation    string
	RequiredFeatures string

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

func (c *Content) ToStore(dest store.Store) error {
	c.Fill()
	if dmeta, ok := dest.(store.MetaSaver); ok {
		meta := map[string]string{
			"Name":             c.Meta.Name,
			"Source":           c.Meta.Source,
			"Description":      c.Meta.Description,
			"Version":          c.Meta.Version,
			"Type":             c.Meta.Type,
			"Documentation":    c.Meta.Documentation,
			"RequiredFeatures": c.Meta.RequiredFeatures,
		}
		if err := dmeta.SetMetaData(meta); err != nil {
			return err
		}
	}
	for section, vals := range c.Sections {
		sub, err := dest.MakeSub(section)
		if err != nil {
			return err
		}
		for k, v := range vals {
			if err := sub.Save(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (c *Content) FromStore(src store.Store) error {
	c.Fill()
	if smeta, ok := src.(store.MetaSaver); ok {
		for k, v := range smeta.MetaData() {
			switch k {
			case "Name":
				c.Meta.Name = v
			case "Source":
				c.Meta.Source = v
			case "Description":
				c.Meta.Description = v
			case "Version":
				c.Meta.Version = v
			case "Type":
				c.Meta.Type = v
			case "Documentation":
				c.Meta.Documentation = v
			case "RequiredFeatures":
				c.Meta.RequiredFeatures = v
			}
		}
	}
	for section, subStore := range src.Subs() {
		if _, err := New(section); err != nil {
			continue
		}
		keys, err := subStore.Keys()
		if err != nil {
			return err
		}
		c.Sections[section] = map[string]interface{}{}
		for _, key := range keys {
			val, _ := New(section)
			if f, ok := val.(Filler); ok {
				f.Fill()
			}
			if err := subStore.Load(key, val); err != nil {
				return err
			}
			c.Sections[section][key] = val
		}
	}

	c.Meta.Type, c.Meta.Overwritable, c.Meta.Writable = getExtraFields(c.Key(), c.Meta.Type)
	return nil
}

type Sections map[string]Section
type Section map[string]interface{}

func (c *Content) Prefix() string {
	return "contents"
}

func (c *Content) Key() string {
	return c.Meta.Name
}

func (c *Content) KeyName() string {
	return "Meta.Name"
}

func (c *Content) Fill() {
	if c.Sections == nil {
		c.Sections = Sections(map[string]Section{})
	}
}

func (c *Content) AuthKey() string {
	return c.Key()
}

// swagger:model
type ContentSummary struct {
	Meta     ContentMetaData `json:"meta"`
	Counts   map[string]int
	Warnings []string
}

func (c *ContentSummary) Fill() {
	if c.Counts == nil {
		c.Counts = map[string]int{}
	}
	if c.Warnings == nil {
		c.Warnings = []string{}
	}
}

func (c *ContentSummary) FromStore(src store.Store) {
	c.Fill()
	if smeta, ok := src.(store.MetaSaver); ok {
		for k, v := range smeta.MetaData() {
			switch k {
			case "Name":
				c.Meta.Name = v
			case "Source":
				c.Meta.Source = v
			case "Description":
				c.Meta.Description = v
			case "Version":
				c.Meta.Version = v
			case "Type":
				c.Meta.Type = v
			case "Documentation":
				c.Meta.Documentation = v
			case "RequiredFeatures":
				c.Meta.RequiredFeatures = v
			}
		}
	}
	for section, subStore := range src.Subs() {
		keys, err := subStore.Keys()
		if err != nil {
			continue
		}
		c.Counts[section] = len(keys)
	}

	c.Meta.Type, c.Meta.Overwritable, c.Meta.Writable = getExtraFields(c.Meta.Name, c.Meta.Type)
	return
}

// Return type, overwritable, writable
func getExtraFields(n, t string) (string, bool, bool) {
	writable := false
	overwritable := false
	if t != "" {
		if t == "default" {
			overwritable = true
		}
	} else {
		t = "dynamic"
	}
	if n == "BackingStore" {
		t = "writable"
		writable = true
	} else if n == "LocalStore" {
		t = "local"
		overwritable = true
	} else if n == "BasicStore" {
		t = "basic"
		overwritable = true
	} else if n == "DefaultStore" {
		t = "default"
		overwritable = true
	}
	return t, overwritable, writable
}
