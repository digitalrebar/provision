package models

import (
	"fmt"
	"strings"

	"github.com/digitalrebar/store"
	"github.com/gofunky/semver"
)

// All fields must be strings
// All string fields will be trimmed except Documentation.
type ContentMetaData struct {
	// required: true
	Name        string
	Version     string // If present, must be parseable as semver
	Description string
	Source      string // Was who authored it, but was confusing

	// Optional fields
	Documentation    string
	RequiredFeatures string

	// New descriptor fields for catalog
	Color         string
	Icon          string
	Author        string
	DisplayName   string
	License       string
	Copyright     string
	CodeSource    string
	Order         string
	Tags          string // Comma separated list
	DocUrl        string
	Prerequisites string // also a comma-seperated list. May contain semver

	// Informational Fields
	Type         string
	Writable     bool
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

func ParseContentPrerequisites(prereqs string) (map[string]semver.Range, error) {
	res := map[string]semver.Range{}
	for _, v := range strings.Split(prereqs, ",") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		parts := strings.SplitN(v, ":", 2)
		if len(parts) == 1 {
			parts = append(parts, ">=0.0.0")
		}
		ver, err := semver.ParseRange(strings.TrimSpace(parts[1]))
		if err != nil {
			return nil, fmt.Errorf("Invalid version requirement for %s: %v", parts[0], err)
		}
		res[strings.TrimSpace(parts[0])] = ver
	}
	return res, nil
}

func (c *Content) ToStore(dest store.Store) error {
	c.Fill()
	if dmeta, ok := dest.(store.MetaSaver); ok {
		meta := map[string]string{
			"Name":        strings.TrimSpace(c.Meta.Name),
			"Version":     strings.TrimSpace(c.Meta.Version),
			"Description": strings.TrimSpace(c.Meta.Description),
			"Source":      strings.TrimSpace(c.Meta.Source),

			"Type": c.Meta.Type,

			"Documentation":    c.Meta.Documentation,
			"RequiredFeatures": strings.TrimSpace(c.Meta.RequiredFeatures),

			"Color":         strings.TrimSpace(c.Meta.Color),
			"Icon":          strings.TrimSpace(c.Meta.Icon),
			"Author":        strings.TrimSpace(c.Meta.Author),
			"DisplayName":   strings.TrimSpace(c.Meta.DisplayName),
			"License":       strings.TrimSpace(c.Meta.License),
			"Copyright":     strings.TrimSpace(c.Meta.Copyright),
			"CodeSource":    strings.TrimSpace(c.Meta.CodeSource),
			"Order":         strings.TrimSpace(c.Meta.Order),
			"Tags":          strings.TrimSpace(c.Meta.Tags),
			"DocUrl":        strings.TrimSpace(c.Meta.DocUrl),
			"Prerequisites": strings.TrimSpace(c.Meta.Prerequisites),
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

func (c *Content) Mangle(thunk func(string, interface{}) (interface{}, error)) error {
	for section := range c.Sections {
		for k := range c.Sections[section] {
			if final, err := thunk(section, c.Sections[section][k]); err == nil && final != nil {
				c.Sections[section][k] = final
			} else if err != nil {
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
			tv := strings.TrimSpace(v)
			switch k {
			case "Name":
				c.Meta.Name = tv
			case "Source":
				c.Meta.Source = tv
			case "Description":
				c.Meta.Description = tv
			case "Version":
				c.Meta.Version = tv
			case "Type":
				c.Meta.Type = tv
			case "Documentation":
				c.Meta.Documentation = v
			case "RequiredFeatures":
				c.Meta.RequiredFeatures = tv
			case "Color":
				c.Meta.Color = tv
			case "Icon":
				c.Meta.Icon = tv
			case "Author":
				c.Meta.Author = tv
			case "DisplayName":
				c.Meta.DisplayName = tv
			case "License":
				c.Meta.License = tv
			case "Copyright":
				c.Meta.Copyright = tv
			case "CodeSource":
				c.Meta.CodeSource = tv
			case "Order":
				c.Meta.Order = tv
			case "Tags":
				c.Meta.Tags = tv
			case "DocUrl":
				c.Meta.DocUrl = tv
			case "Prerequisites":
				c.Meta.Prerequisites = tv
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
			tv := strings.TrimSpace(v)
			switch k {
			case "Name":
				c.Meta.Name = tv
			case "Source":
				c.Meta.Source = tv
			case "Description":
				c.Meta.Description = tv
			case "Version":
				c.Meta.Version = tv
			case "Type":
				c.Meta.Type = tv
			case "Documentation":
				c.Meta.Documentation = v
			case "RequiredFeatures":
				c.Meta.RequiredFeatures = tv
			case "Color":
				c.Meta.Color = tv
			case "Icon":
				c.Meta.Icon = tv
			case "Author":
				c.Meta.Author = tv
			case "DisplayName":
				c.Meta.DisplayName = tv
			case "License":
				c.Meta.License = tv
			case "Copyright":
				c.Meta.Copyright = tv
			case "CodeSource":
				c.Meta.CodeSource = tv
			case "Order":
				c.Meta.Order = tv
			case "DocUrl":
				c.Meta.DocUrl = tv
			case "Prerequisites":
				c.Meta.Prerequisites = tv
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
