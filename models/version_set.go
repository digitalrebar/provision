package models

import (
	"bytes"
	"encoding/json"
	"log"
)

/*
 * VersionSet tracks a versioned thing in the RackN system
 */

type FileData struct {
	Path      string // Including name e.g. files/raid/jjj
	Sha256Sum string // sha256sum of the item
	Source    string // Location to get file from - self://path is this DRP
	Explode   bool
}

// VersionSet structure that handles RawModel instead of dealing with
// RawModel which is how DRP is storing it.
//
// An element with Version = ignore means leave it loaded.
type VersionSet struct {
	Validation
	Access
	Meta
	Owned
	Bundled

	Id string `index:",key"`

	Documentation string
	Description   string

	Apply        bool
	DRPVersion   string
	DRPUXVersion string
	Components   []*Element
	Plugins      []*Plugin
	Prefs        map[string]string
	Files        []*FileData
	Global       map[string]interface{}
}

func (vs *VersionSet) Key() string {
	return vs.Id
}

func (vs *VersionSet) KeyName() string {
	return "Id"
}

func (vs *VersionSet) AuthKey() string {
	return vs.Key()
}

func (vs *VersionSet) Prefix() string {
	return "version_sets"
}

// GetDocumentaiton returns the object's Documentation
func (vs *VersionSet) GetDocumentation() string {
	return vs.Documentation
}

// GetDescription returns the object's Description
func (vs *VersionSet) GetDescription() string {
	return vs.Description
}

// Clone the VersionSet
func (vs *VersionSet) Clone() *VersionSet {
	ci2 := &VersionSet{}
	buf := bytes.Buffer{}
	enc, dec := json.NewEncoder(&buf), json.NewDecoder(&buf)
	if err := enc.Encode(vs); err != nil {
		log.Panicf("Failed to encode endpoint:%s: %v", vs.Id, err)
	}
	if err := dec.Decode(ci2); err != nil {
		log.Panicf("Failed to decode endpoint:%s: %v", vs.Id, err)
	}
	return ci2
}

func (vs *VersionSet) Fill() {
	vs.Validation.fill(vs)
	if vs.Meta == nil {
		vs.Meta = Meta{}
	}
	if vs.Errors == nil {
		vs.Errors = []string{}
	}
	if vs.Components == nil {
		vs.Components = []*Element{}
	}
	if vs.Plugins == nil {
		vs.Plugins = []*Plugin{}
	}
	if vs.Prefs == nil {
		vs.Prefs = map[string]string{}
	}
	if vs.Global == nil {
		vs.Global = map[string]interface{}{}
	}
	if vs.Files == nil {
		vs.Files = []*FileData{}
	}
}

func (vs *VersionSet) Merge(nvs *VersionSet) {
	vs.Apply = vs.Apply && nvs.Apply
	if nvs.DRPVersion != "" {
		vs.DRPVersion = nvs.DRPVersion
	}
	if nvs.DRPUXVersion != "" {
		vs.DRPUXVersion = nvs.DRPUXVersion
	}
	// Add in components
	for _, c := range nvs.Components {
		found := false
		for i, curc := range vs.Components {
			if curc.Name == c.Name {
				vs.Components[i] = c
				found = true
				break
			}
		}
		if !found {
			vs.Components = append(vs.Components, c)
		}
	}

	// Add in plugins
	for _, c := range nvs.Plugins {
		found := false
		for i, curc := range vs.Plugins {
			if curc.Name == c.Name {
				vs.Plugins[i] = c
				found = true
				break
			}
		}
		if !found {
			vs.Plugins = append(vs.Plugins, c)
		}
	}

	// Merge prefs
	for k, v := range nvs.Prefs {
		vs.Prefs[k] = v
	}

	// Merge Global
	for k, v := range nvs.Global {
		vs.Global[k] = v
	}

	// Files
	for _, c := range nvs.Files {
		found := false
		for i, curc := range vs.Files {
			if curc.Path == c.Path {
				vs.Files[i] = c
				found = true
				break
			}
		}
		if !found {
			vs.Files = append(vs.Files, c)
		}
	}
}

func (vs *VersionSet) SliceOf() interface{} {
	s := []*VersionSet{}
	return &s
}

func (vs *VersionSet) ToModels(obj interface{}) []Model {
	items := obj.(*[]*VersionSet)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (vs *VersionSet) CanHaveActions() bool {
	return true
}

// SetName sets the name. In this case, it sets Id.
func (vs *VersionSet) SetName(name string) {
	vs.Id = name
}
