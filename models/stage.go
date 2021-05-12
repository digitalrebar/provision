package models

import "strings"

// Stage encapsulates a set of tasks and profiles to apply
// to a Machine in a BootEnv.
//
// swagger:model
type Stage struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	Partialed
	// The name of the stage.
	//
	// required: true
	Name string `index:",key"`
	// A description of this stage.  This should tell what it is for,
	// any special considerations that should be taken into account when
	// using it, etc.
	Description string
	// Documentation of this stage.  This should tell what
	// the stage is for, any special considerations that
	// should be taken into account when using it, etc. in rich structured text (rst).
	Documentation string
	// The templates that should be expanded into files for the stage.
	//
	// required: true
	Templates []TemplateInfo
	// The list of extra required parameters for this
	// stage. They should be present as Machine.Params when
	// the stage is applied to the machine.
	//
	// required: true
	RequiredParams []string
	// The list of extra optional parameters for this
	// stage. They can be present as Machine.Params when
	// the stage is applied to the machine.  These are more
	// other consumers of the stage to know what parameters
	// could additionally be applied to the stage by the
	// renderer based upon the Machine.Params
	//
	OptionalParams []string
	// Params contains parameters for the stage.
	// This allows the machine to access these values while in this stage.
	Params map[string]interface{}
	// The BootEnv the machine should be in to run this stage.
	// If the machine is not in this bootenv, the bootenv of the
	// machine will be changed.
	//
	// required: true
	BootEnv string
	// The list of initial machine tasks that the stage should run
	Tasks []string
	// The list of profiles a machine should use while in this stage.
	// These are used after machine profiles, but before global.
	Profiles []string
	// Flag to indicate if a node should be PXE booted on this
	// transition into this Stage.  The nextbootpxe and reboot
	// machine actions will be called if present and Reboot is true
	Reboot bool
	// This flag is deprecated and will always be TRUE.
	RunnerWait bool
}

// GetMeta gets the meta data for this object
func (s *Stage) GetMeta() Meta {
	return s.Meta
}

// SetMeta sets the meta data for this object
func (s *Stage) SetMeta(d Meta) {
	s.Meta = d
}

// GetDocumentaiton returns the object's Documentation
func (s *Stage) GetDocumentation() string {
	return s.Documentation
}

// GetDescription returns the object's Description
func (s *Stage) GetDescription() string {
	return s.Description
}

// Validate makes sure that the object is valid (outside of references)
func (s *Stage) Validate() {
	s.AddError(ValidName("Invalid Name", s.Name))
	if s.BootEnv != "" {
		s.AddError(ValidName("Invalid BootEnv", s.BootEnv))
	}

	for _, p := range s.RequiredParams {
		s.AddError(ValidParamName("Invalid Required Param", p))
	}
	for _, p := range s.OptionalParams {
		s.AddError(ValidParamName("Invalid Optional Param", p))
	}
	tmplNames := map[string]int{}
	for i := range s.Templates {
		tmpl := &(s.Templates[i])
		tmpl.SanityCheck(i, s, false)
		if j, ok := tmplNames[tmpl.Name]; ok {
			s.Errorf("Template %d and %d have the same name %s", i, j, tmpl.Name)
		} else {
			tmplNames[tmpl.Name] = i
		}
	}
	for _, p := range s.Profiles {
		s.AddError(ValidName("Invalid Profile", p))
	}
	for _, t := range s.Tasks {
		if parts := strings.SplitN(t, ":", 2); len(parts) == 1 {
			s.AddError(ValidName("Invalid Task", t))
		} else {
			switch parts[0] {
			case "action":
				pparts := strings.SplitN(parts[1], ":", 2)
				if len(pparts) != 2 {
					s.Errorf("Invalid action specifier %s", parts[1])
					continue
				}
				s.AddError(ValidName("Invalid Plugin", pparts[0]))
				s.AddError(ValidName("Invalid Action", pparts[1]))
			case "chroot", "bootenv", "stage", "context":
			default:
				s.Errorf("Invalid Task: %s", t)
			}
		}
	}
}

// Prefix returns the object type
func (s *Stage) Prefix() string {
	return "stages"
}

// Key returns the primary index for this object
func (s *Stage) Key() string {
	return s.Name
}

// KeyName returns the field of the object that is used as the primary key
func (s *Stage) KeyName() string {
	return "Name"
}

// Fill initializes the object
func (s *Stage) Fill() {
	s.Validation.fill(s)
	if s.Meta == nil {
		s.Meta = Meta{}
	}
	if s.Templates == nil {
		s.Templates = []TemplateInfo{}
	}
	if s.RequiredParams == nil {
		s.RequiredParams = []string{}
	}
	if s.OptionalParams == nil {
		s.OptionalParams = []string{}
	}
	if s.Tasks == nil {
		s.Tasks = []string{}
	}
	if s.Profiles == nil {
		s.Profiles = []string{}
	}
	if s.Params == nil {
		s.Params = map[string]interface{}{}
	}
}

// AuthKey returns the value that should be validated against claims
func (s *Stage) AuthKey() string {
	return s.Key()
}

// SliceOf returns an empty slice of this type of objects
func (s *Stage) SliceOf() interface{} {
	s2 := []*Stage{}
	return &s2
}

// ToModels converts a slice of these specific objects to a slice of Model interfaces
func (s *Stage) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Stage)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

// match Profiler interface

// GetProfiles gets the profiles on this stage
func (s *Stage) GetProfiles() []string {
	return s.Profiles
}

// SetProfiles sets the profiles on this stage
func (s *Stage) SetProfiles(p []string) {
	s.Profiles = p
}

// match Paramer interface

// GetParams gets the parameters on this stage
func (s *Stage) GetParams() map[string]interface{} {
	return copyMap(s.Params)
}

// SetParams sets the parameters on this stage
func (s *Stage) SetParams(p map[string]interface{}) {
	s.Params = copyMap(p)
}

// match BootEnver interface

// GetBootEnv gets the name of the bootenv on this stage
func (s *Stage) GetBootEnv() string {
	return s.BootEnv
}

// SetBootEnv sets the bootenv on the stage
func (s *Stage) SetBootEnv(be string) {
	s.BootEnv = be
}

// match TaskRunner interface

// GetTasks returns the tasks associated with this stage
func (s *Stage) GetTasks() []string {
	return s.Tasks
}

// SetTasks sets the tasks in this stage
func (s *Stage) SetTasks(t []string) {
	s.Tasks = t
}

// SetName sets the name of the object
func (s *Stage) SetName(n string) {
	s.Name = n
}

// CanHaveActions returns if the object can have actions from plugins
func (s *Stage) CanHaveActions() bool {
	return true
}
