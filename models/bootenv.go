package models

// OsInfo holds information about the operating system this BootEnv
// maps to.  Most of this information is optional for now.
// swagger:model
type OsInfo struct {
	// The name of the OS this BootEnv has.
	//
	// required: true
	Name string
	// The family of operating system (linux distro lineage, etc)
	Family string
	// The codename of the OS, if any.
	Codename string
	// The version of the OS, if any.
	Version string
	// The name of the ISO that the OS should install from.
	IsoFile string
	// The SHA256 of the ISO file.  Used to check for corrupt downloads.
	IsoSha256 string
	// The URL that the ISO can be downloaded from, if any.
	//
	// swagger:strfmt uri
	IsoUrl string
}

// BootEnv encapsulates the machine-agnostic information needed by the
// provisioner to set up a boot environment.
//
// swagger:model
type BootEnv struct {
	Validation
	// The name of the boot environment.  Boot environments that install
	// an operating system must end in '-install'.
	//
	// required: true
	Name string
	// A description of this boot environment.  This should tell what
	// the boot environment is for, any special considerations that
	// shoudl be taken into account when using it, etc.
	Description string
	// The OS specific information for the boot environment.
	OS OsInfo
	// The templates that should be expanded into files for the
	// boot environment.
	//
	// required: true
	Templates []TemplateInfo
	// The partial path to the kernel for the boot environment.  This
	// should be path that the kernel is located at in the OS ISO or
	// install archive.
	//
	// required: true
	Kernel string
	// Partial paths to the initrds that should be loaded for the boot
	// environment. These should be paths that the initrds are located
	// at in the OS ISO or install archive.
	//
	// required: true
	Initrds []string
	// A template that will be expanded to create the full list of
	// boot parameters for the environment.
	//
	// required: true
	BootParams string
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
	// OnlyUnknown indicates whether this bootenv can be used without a
	// machine.  Only bootenvs with this flag set to `true` be used for
	// the unknownBootEnv preference.
	//
	// required: true
	OnlyUnknown bool
	// The list of initial machine tasks that the boot environment should get
	Tasks []string
}

func (b *BootEnv) Prefix() string {
	return "bootenvs"
}

func (b *BootEnv) Key() string {
	return b.Name
}

func (b *BootEnv) AuthKey() string {
	return b.Key()
}
