package models

import (
	"strconv"
	"strings"
)

// OsInfo holds information about the operating system this BootEnv
// maps to.  Most of this information is optional for now.
// swagger:model
type OsInfo struct {
	// The name of the OS this BootEnv has.  It should be formatted as
	// family-version.
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

// FamilyName is a helper that figures out the family (read: distro
// name) of the OS.  It uses Family if set, the first part of the Name
// otherwise.
func (o OsInfo) FamilyName() string {
	if o.Family != "" {
		return o.Family
	}
	return strings.Split(o.Name, "-")[0]
}

// FamilyType figures out the lineage of the OS.  If the OS is
// descended from RHEL, then "rhel" is returned.  If the OS is
// descended from Debian, then "debian" is returned, otherwise
// FamilyName() is returned.  Return values of this function are
// subject to change as support for new distros is brought onboard.
func (o OsInfo) FamilyType() string {
	switch o.FamilyName() {
	case "centos", "redhat", "fedora", "scientificlinux":
		return "rhel"
	case "debian", "ubuntu":
		return "debian"
	default:
		return o.FamilyName()
	}
}

// FamilyVersion figures out the version of the OS.  It returns the
// Version field if set, and the second part of the OS name if not
// set.  THis should be a Semver-ish version string, not a codename,
// release name, or similar item.
func (o OsInfo) FamilyVersion() string {
	if o.Version != "" {
		return o.Version
	}
	parts := strings.Split(o.Name, "-")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

// VersionEq returns true of this OS version is equal to the degree of
// accuracy implied by other -- o.Version(7.3) is VersionEq to 7 and
// 7.3, but not 7.3.11
func (o OsInfo) VersionEq(other string) bool {
	partCmp := func(a, b string) bool {
		ai, aerr := strconv.ParseInt(a, 10, 64)
		bi, berr := strconv.ParseInt(b, 10, 64)
		if aerr == nil && berr == nil {
			return ai == bi
		}
		return a == b
	}
	myParts := strings.Split(o.FamilyVersion(), ".")
	otherParts := strings.Split(other, ".")
	if len(myParts) < len(otherParts) {
		return false
	}

	for i := 0; i < len(otherParts); i++ {
		if !partCmp(myParts[i], otherParts[i]) {
			return false
		}
	}
	return true
}

// BootEnv encapsulates the machine-agnostic information needed by the
// provisioner to set up a boot environment.
//
// swagger:model
type BootEnv struct {
	Validation
	Access
	Meta
	// The name of the boot environment.  Boot environments that install
	// an operating system must end in '-install'.
	//
	// required: true
	Name string
	// A description of this boot environment.  This should tell what
	// the boot environment is for, any special considerations that
	// should be taken into account when using it, etc.
	Description string
	// Documentation of this boot environment.  This should tell what
	// the boot environment is for, any special considerations that
	// should be taken into account when using it, etc. in rich structured text (rst).
	Documentation string
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
}

func (b *BootEnv) GetMeta() Meta {
	return b.Meta
}

func (b *BootEnv) SetMeta(d Meta) {
	b.Meta = d
}

func (b *BootEnv) GetDocumentation() string {
	return b.Documentation
}

func (b *BootEnv) Validate() {
	b.AddError(ValidName("Invalid Name", b.Name))
	for _, p := range b.RequiredParams {
		b.AddError(ValidParamName("Invalid Required Param", p))
	}
	for _, p := range b.OptionalParams {
		b.AddError(ValidParamName("Invalid Optional Param", p))
	}
	tmplNames := map[string]int{}
	for i := range b.Templates {
		tmpl := &(b.Templates[i])
		tmpl.SanityCheck(i, b, false)
		if j, ok := tmplNames[tmpl.Name]; ok {
			b.Errorf("Template %d and %d have the same name %s", i, j, tmpl.Name)
		} else {
			tmplNames[tmpl.Name] = i
		}
	}
}

func (b *BootEnv) Prefix() string {
	return "bootenvs"
}

func (b *BootEnv) Key() string {
	return b.Name
}

func (b *BootEnv) KeyName() string {
	return "Name"
}

func (b *BootEnv) AuthKey() string {
	return b.Key()
}

func (b *BootEnv) SliceOf() interface{} {
	s := []*BootEnv{}
	return &s
}

func (b *BootEnv) ToModels(obj interface{}) []Model {
	items := obj.(*[]*BootEnv)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

func (b *BootEnv) Fill() {
	b.Validation.fill()
	if b.Meta == nil {
		b.Meta = Meta{}
	}
	if b.Initrds == nil {
		b.Initrds = []string{}
	}
	if b.OptionalParams == nil {
		b.OptionalParams = []string{}
	}
	if b.RequiredParams == nil {
		b.RequiredParams = []string{}
	}
	if b.Templates == nil {
		b.Templates = []TemplateInfo{}
	}
}

func (b *BootEnv) SetName(n string) {
	b.Name = n
}

func (b *BootEnv) CanHaveActions() bool {
	return true
}

func (b *BootEnv) NetBoot() bool {
	return b.OnlyUnknown || b.Kernel != ""
}
