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

func (s *Stage) Prefix() string {
	return "stages"
}

func (s *Stage) Key() string {
	return s.Name
}

func (s *Stage) AuthKey() string {
	return s.Key()
}
