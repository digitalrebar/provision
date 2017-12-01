package models

// Stage encapsulates a set of tasks and profiles to apply
// to a Machine in a BootEnv.
//
// swagger:model
type Stage struct {
	Validation
	Access
	MetaData
	// The name of the boot environment.  Boot environments that install
	// an operating system must end in '-install'.
	//
	// required: true
	Name string
	// A description of this boot environment.  This should tell what
	// the boot environment is for, any special considerations that
	// shoudl be taken into account when using it, etc.
	Description string
	// The templates that should be expanded into files for the stage.
	//
	// required: true
	Templates []TemplateInfo
	// The list of extra required parameters for this
	// bootstate. They should be present as Machine.Params when
	// the bootenv is applied to the machine.
	//
	// required: true
	RequiredParams []string
	// The list of extra optional parameters for this
	// bootstate. They can be present as Machine.Params when
	// the bootenv is applied to the machine.  These are more
	// other consumers of the bootenv to know what parameters
	// could additionally be applied to the bootenv by the
	// renderer based upon the Machine.Params
	//
	OptionalParams []string
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
	// Flag to indicate if the runner should wait for more tasks
	// while in this stage.
	RunnerWait bool
}

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
	for _, t := range s.Templates {
		s.AddError(ValidName("Invalid Template Name", t.Name))
	}
	for _, p := range s.Profiles {
		s.AddError(ValidName("Invalid Profile", p))
	}
	for _, t := range s.Tasks {
		s.AddError(ValidName("Invalid Task", t))
	}
}

func (s *Stage) Prefix() string {
	return "stages"
}

func (s *Stage) Key() string {
	return s.Name
}

func (s *Stage) Fill() {
	s.Validation.fill()
	s.MetaData.fill()
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
}

func (s *Stage) AuthKey() string {
	return s.Key()
}

func (b *Stage) SliceOf() interface{} {
	s := []*Stage{}
	return &s
}

func (b *Stage) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Stage)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

// match Profiler interface
func (b *Stage) GetProfiles() []string {
	return b.Profiles
}

func (b *Stage) SetProfiles(p []string) {
	b.Profiles = p
}

// match BootEnver interface
func (b *Stage) GetBootEnv() string {
	return b.BootEnv
}

func (b *Stage) SetBootEnv(s string) {
	b.BootEnv = s
}

// match TaskRunner interface
func (b *Stage) GetTasks() []string {
	return b.Tasks
}

func (b *Stage) SetTasks(t []string) {
	b.Tasks = t
}

func (b *Stage) SetName(n string) {
	b.Name = n
}
