package models

import (
	"sort"
	"strings"
)

// Task is a thing that can run on a Machine.
//
// swagger:model
type Task struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	// Name is the name of this Task.  Task names must be globally unique
	//
	// required: true
	Name string `index:",key"`
	// Description is a one-line description of this Task.
	Description string
	// Documentation should describe in detail what this task should do on a machine.
	Documentation string
	// Templates lists the templates that need to be rendered for the Task.
	//
	// required: true
	Templates []TemplateInfo
	// RequiredParams is the list of parameters that are required to be present on
	// Machine.Params or in a profile attached to the machine.
	//
	// required: true
	RequiredParams []string
	// OptionalParams are extra optional parameters that a template rendered for
	// the Task may use.
	//
	// required: true
	OptionalParams []string
	// Prerequisites are tasks that must have been run in the current
	// BootEnv before this task can be run.
	Prerequisites []string
	// ExtraRoles is a list of Roles whose Claims should be added to the default
	// set of allowable Claims when a Job based on this task is running.
	ExtraRoles []string `json:"ExtraRoles,omitempty"`
	// ExtraClaims is a raw list of Claims that should be added to the default
	// set of allowable Claims when a Job based on this task is running.
	// Any extra claims added here will be added _after_ any added by ExtraRoles
	ExtraClaims []*Claim `json:"ExtraClaims,omitempty"`
}

var (
	validGOOS = map[string]struct{}{
		"any":     struct{}{},
		"darwin":  struct{}{},
		"freebsd": struct{}{},
		"linux":   struct{}{},
		"netbsd":  struct{}{},
		"openbsd": struct{}{},
		"solaris": struct{}{},
		"windows": struct{}{},
		"esxi":    struct{}{}, // Not really GOOS, but a valid template
	}

	validGOARCH = map[string]struct{}{
		"386":     struct{}{},
		"amd64":   struct{}{},
		"arm":     struct{}{},
		"arm64":   struct{}{},
		"ppc64le": struct{}{},
	}
)

func (t *Task) GetMeta() Meta {
	return t.Meta
}

func (t *Task) SetMeta(d Meta) {
	t.Meta = d
}

// GetDocumentaiton returns the object's Documentation
func (t *Task) GetDocumentation() string {
	return t.Documentation
}

// GetDescription returns the object's Description
func (t *Task) GetDescription() string {
	return t.Description
}

func (t *Task) Validate() {
	t.AddError(ValidName("Invalid Name", t.Name))

	for _, p := range t.RequiredParams {
		t.AddError(ValidParamName("Invalid Required Param", p))
	}
	for _, p := range t.OptionalParams {
		t.AddError(ValidParamName("Invalid Optional Param", p))
	}
	printedValidValues := false
	osMetaCount := 0
	tmplNames := map[string]int{}
	for i := range t.Templates {
		tmpl := &(t.Templates[i])
		tmpl.SanityCheck(i, t, true)
		if j, ok := tmplNames[tmpl.Name]; ok {
			t.Errorf("Template %d and %d have the same name %s", i, j, tmpl.Name)
		} else {
			tmplNames[tmpl.Name] = i
		}
		if _, ok := tmpl.Meta["OS"]; ok {
			osMetaCount++
			oses := strings.Split(tmpl.Meta["OS"], ",")
			for _, os := range oses {
				if _, ok := validGOOS[strings.ToLower(strings.TrimSpace(os))]; ok {
					continue
				}
				t.Errorf("Template[%d]: Invalid OS value %s", i, os)
				if !printedValidValues {
					validOSes := make([]string, 0, len(validGOOS))
					for k := range validGOOS {
						validOSes = append(validOSes, k)
					}
					sort.Strings(validOSes)
					t.Errorf("Valid values are: %s", strings.Join(validOSes, ","))
					printedValidValues = true
				}
			}
		}
	}
	if osMetaCount != 0 && osMetaCount != len(tmplNames) {
		t.Errorf("Cannot mix templates with OS metadata and templates without OS metadata")
		for i := range t.Templates {
			tmpl := &(t.Templates[i])
			if _, ok := tmpl.Meta["OS"]; ok {
				t.Errorf("Template[%d] %s has OS metadata %s", i, tmpl.Name, tmpl.Meta["OS"])
			} else {
				t.Errorf("Template[%d] %s is missing OS metadata", i, tmpl.Name)
			}
		}
	}
}

func (t *Task) Prefix() string {
	return "tasks"
}

func (t *Task) Key() string {
	return t.Name
}

func (t *Task) KeyName() string {
	return "Name"
}

func (t *Task) Fill() {
	t.Validation.fill(t)
	if t.Meta == nil {
		t.Meta = Meta{}
	}
	if t.Templates == nil {
		t.Templates = []TemplateInfo{}
	}
	if t.RequiredParams == nil {
		t.RequiredParams = []string{}
	}
	if t.OptionalParams == nil {
		t.OptionalParams = []string{}
	}
	if t.Prerequisites == nil {
		t.Prerequisites = []string{}
	}
	if t.ExtraRoles == nil {
		t.ExtraRoles = []string{}
	}
	if t.ExtraClaims == nil {
		t.ExtraClaims = []*Claim{}
	}
}

func (t *Task) AuthKey() string {
	return t.Key()
}

func (t *Task) SliceOf() interface{} {
	s := []*Task{}
	return &s
}

func (t *Task) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Task)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (t *Task) SetName(n string) {
	t.Name = n
}

func (t *Task) CanHaveActions() bool {
	return true
}
