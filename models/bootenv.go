package models

import (
	"strconv"
	"strings"
)

// ArchInfo tracks information required to make a BootEnv work across
// different system architectures.  It supersedes the matching fields
// in the BootEnv struct and the OsInfo struct.
type ArchInfo struct {
	// IsoFile is the name of the ISO file (or other archive)
	// that contains all the necessary information to be able to
	// boot into this BootEnv for a given arch.
	// At a minimum, it must contain a kernel and initrd that
	// can be booted over the network.
	IsoFile string
	// Sha256 should contain the SHA256 checksum for the IsoFile.
	// If it does, the IsoFile will be checked upon upload to make sure
	// it has not been corrupted.
	Sha256 string
	// IsoUrl is the location that IsoFile can be downloaded from, if any.
	// This must be a full URL, including the filename.
	//
	// swagger:strfmt url
	IsoUrl string
	// The partial path to the kernel for the boot environment.  This
	// should be path that the kernel is located at in the OS ISO or
	// install archive.  If empty, this will fall back to the top-level
	// Kernel field in the BootEnv
	//
	// required: true
	Kernel string
	// Partial paths to the initrds that should be loaded for the boot
	// environment. These should be paths that the initrds are located
	// at in the OS ISO or install archive.  If empty, this will fall back
	// to the top-level Initrds field in the BootEnv
	//
	// required: true
	Initrds []string
	// A template that will be expanded to create the full list of
	// boot parameters for the environment.  If empty, this will fall back
	// to the top-level BootParams field in the BootEnv
	//
	// required: true
	BootParams string
	// Loader is the bootloader that should be used for this boot
	// environment.  If left unspecified and not overridden by a subnet
	// or reservation option, the following boot loaders will be used:
	//
	// * lpxelinux.0 on 386-pcbios platforms that are not otherwise using ipxe.
	//
	// * ipxe.pxe on 386-pcbios platforms that already use ipxe.
	//
	// * ipxe.efi on amd64 EFI platforms.
	//
	// * ipxe-arm64.efi on arm64 EFI platforms.
	//
	// This setting will be overridden by Subnet and Reservation
	// options, and it will also only be in effect when dr-provision is
	// the DHCP server of record.
	Loader string
}

func (a *ArchInfo) Fill() {
	if a.Initrds == nil {
		a.Initrds = []string{}
	}
}

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
	// The name of the ISO that the OS should install from.  If
	// non-empty, this is assumed to be for the amd64 hardware
	// architecture.
	IsoFile string
	// The SHA256 of the ISO file.  Used to check for corrupt downloads.
	// If non-empty, this is assumed to be for the amd64 hardware
	// architecture.
	IsoSha256 string
	// The URL that the ISO can be downloaded from, if any.  If
	// non-empty, this is assumed to be for the amd64 hardware
	// architecture.
	//
	// swagger:strfmt uri
	IsoUrl string
	// SupportedArchitectures maps from hardware architecture (named
	// according to the distro architecture naming scheme) to the
	// architecture-specific parameters for this OS.  If
	// SupportedArchitectures is left empty, then the system assumes
	// that the BootEnv only supports amd64 platforms.
	SupportedArchitectures map[string]ArchInfo
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
// set.  This should be a Semver-ish version string, not a codename,
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
	Owned
	Bundled
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
	// install archive.  Kernel must be non-empty for a BootEnv to be
	// considered net bootable.
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
	// Loaders contains the boot loaders that should be used for various different network
	// boot scenarios.  It consists of a map of machine type -> partial paths to the bootloaders.
	// Valid machine types are:
	//
	// - 386-pcbios for x86 devices using the legacy bios.
	//
	// - amd64-uefi for x86 devices operating in UEFI mode
	//
	// - arm64-uefi for arm64 devices operating in UEFI mode
	//
	// Other machine types will be added as dr-provision gains support for them.
	//
	// If this map does not contain an entry for the machine type, the DHCP server will fall back to
	// the following entries in this order:
	//
	// - The Loader specified in the ArchInfo struct from this BootEnv, if it exists.
	//
	// - The value specified in the bootloaders param for the machine type specified on the machine, if it exists.
	//
	// - The value specified in the bootloaders param in the global profile, if it exists.
	//
	// - The value specified in the default value for the bootloaders param.
	//
	// - One of the following vaiues:
	//
	//   - lpxelinux.0 for 386-pcbios
	//
	//   - ipxe.efi for amd64-uefi
	//
	//   - ipxe-arm64.efi for arm64-uefi
	//
	// required: true
	Loaders map[string]string
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

func (b *BootEnv) GetDescription() string {
	return b.Description
}

// IsoFor is a helper function used by the backend to locate the ISO
// file that should be expanded to provide the install tree required
// for the bootenv to function.
func (b *BootEnv) IsoFor(arch string) string {
	info, ok := b.OS.SupportedArchitectures[arch]
	if ok {
		return info.IsoFile
	}
	if a, _ := SupportedArch(arch); a == "amd64" {
		return b.OS.IsoFile
	}
	return ""
}

// ShaFor is a helper to return the right SHA256 sum for the ISO that
// provides files for the BootEnv.
func (b *BootEnv) ShaFor(arch string) string {
	info, ok := b.OS.SupportedArchitectures[arch]
	if ok {
		return info.Sha256
	}
	if a, _ := SupportedArch(arch); a == "amd64" {
		return b.OS.IsoSha256
	}
	return ""
}

// IsoUrlFor is a helper to return the upstream URL that the ISO for
// the BootEnv can be downloaded from.  This generally points to a
// mirror location on the public Internet if one exists.
func (b *BootEnv) IsoUrlFor(arch string) string {
	info, ok := b.OS.SupportedArchitectures[arch]
	if ok {
		return info.IsoUrl
	}
	if a, _ := SupportedArch(arch); a == "amd64" {
		return b.OS.IsoUrl
	}
	return ""
}

func (b *BootEnv) KernelFor(arch string) string {
	info, ok := b.OS.SupportedArchitectures[arch]
	if ok && info.Kernel != "" {
		return info.Kernel
	}
	return b.Kernel
}

func (b *BootEnv) InitrdsFor(arch string) []string {
	info, ok := b.OS.SupportedArchitectures[arch]
	if ok && len(info.Initrds) > 0 {
		return info.Initrds
	}
	return b.Initrds
}

func (b *BootEnv) BootParamsFor(arch string) string {
	info, ok := b.OS.SupportedArchitectures[arch]
	if ok && info.BootParams != "" {
		return info.BootParams
	}
	return b.BootParams
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
	for k := range b.Loaders {
		switch k {
		case "386-pcbios", "amd64-uefi", "arm64-uefi":
		default:
			b.Errorf("%s is not a supported loader type", k)
		}
	}
	for k := range b.OS.SupportedArchitectures {
		if _, ok := SupportedArch(k); !ok {
			b.Errorf("%s is not a supported architecture", k)
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
	if b.OS.SupportedArchitectures == nil {
		b.OS.SupportedArchitectures = map[string]ArchInfo{}
	}
	if b.Loaders == nil {
		b.Loaders = map[string]string{}
	}
	for k, v := range b.OS.SupportedArchitectures {
		v.Fill()
		b.OS.SupportedArchitectures[k] = v
	}
}

func (b *BootEnv) SetName(n string) {
	b.Name = n
}

func (b *BootEnv) CanHaveActions() bool {
	return true
}

// NetBoot returns whether this bootenv is able to boot via PXE or
// some other network mechanism.
func (b *BootEnv) NetBoot() bool {
	return b.OnlyUnknown || b.Kernel != ""
}
