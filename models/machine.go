package models

import (
	"net"
	"reflect"
	"strings"

	"github.com/pborman/uuid"
)

// SupportedArch normalizes system architectures and returns whether
// it is one we know how to normalize.
func SupportedArch(s string) (string, bool) {
	switch strings.ToLower(s) {
	// rpi4 is a hack because it is really arm64, but reports as amd64.
	// We build a lying bootenv for this purpose.
	case "rpi4":
		return "rpi4", true
	case "amd64", "x86_64":
		return "amd64", true
	case "386", "486", "686", "i386", "i486", "i686":
		return "386", true
	case "arm", "armel", "armhfp":
		return "arm", true
	case "arm64", "aarch64":
		return "arm64", true
	case "ppc64", "power9":
		return "ppc64", true
	case "ppc64le":
		return "ppc64le", true
	case "mips64":
		return "mips64", true
	case "mips64le", "mips64el":
		return "mips64le", true
	case "s390x":
		return "s390x", true
	case "mips":
		return "mips", true
	case "mipsle", "mipsel":
		return "mipsle", true
	default:
		return "", false
	}
}

// ArchEqual returns whether two arches are equal.
func ArchEqual(a, b string) bool {
	a1, aok := SupportedArch(a)
	b1, bok := SupportedArch(b)
	return aok && bok && a1 == b1
}

type MachineFingerprint struct {
	// DMI.System.Manufacturer + DMI.System.ProductName + DMI.System.SerialNumber, SHA256 hashed
	// Hash must not be zero-length to match. 25 points
	SSNHash []byte
	// DMI.System.Manufacturer + DMI.System.ProductName + DMI.Chassis[0].SerialNumber, SHA256 hashed
	// Hash must not be zero-length to match. 25 points
	CSNHash []byte
	// DMI.System.UUID, not hashed. Must be non zero length and must be a non-zero UUID. 50 point match
	SystemUUID string
	// MemoryIds is an array of SHA256sums if the following fields in each
	// entry of the DMI.Memory.Devices array concatenated together:
	//  * Manufacturer
	//  * PartNumber
	//  * SerialNumber
	// Each hash must not be zero length
	// Score is % matched.
	MemoryIds [][]byte
}

// TaskStack contains an task stack for a machine.
// This is used by the error handling code pieces.
type TaskStack struct {
	CurrentTask int
	TaskList    []string
}

func (m *MachineFingerprint) Fill() {
	if m.SSNHash == nil {
		m.SSNHash = []byte{}
	}
	if m.CSNHash == nil {
		m.CSNHash = []byte{}
	}
	if m.MemoryIds == nil {
		m.MemoryIds = [][]byte{}
	}
}

// Machine represents a single bare-metal system that the provisioner
// should manage the boot environment for.
// swagger:model
type Machine struct {
	Validation
	Access
	Meta
	Owned
	Bundled
	Partialed
	// The name of the machine.  This must be unique across all
	// machines, and by convention it is the FQDN of the machine,
	// although nothing enforces that.
	//
	// required: true
	// swagger:strfmt hostname
	Name string `index:",uniq"`
	// A description of this machine.  This can contain any reference
	// information for humans you want associated with the machine.
	Description string
	// The UUID of the machine.
	// This is auto-created at Create time, and cannot change afterwards.
	//
	// required: true
	// swagger:strfmt uuid
	Uuid uuid.UUID `index:",key"`
	// The UUID of the job that is currently running on the machine.
	//
	// swagger:strfmt uuid
	CurrentJob uuid.UUID
	// The IPv4 address of the machine that should be used for PXE
	// purposes.  Note that this field does not directly tie into DHCP
	// leases or reservations -- the provisioner relies solely on this
	// address when determining what to render for a specific machine.
	// Address is updated automatically by the DHCP system if
	// HardwareAddrs is filled out.
	//
	// swagger:strfmt ipv4
	Address net.IP
	// The stage that the Machine is currently in.  If Workflow is also
	// set, this field is read-only, otherwise changing it will change
	// the Stage the system is in.
	Stage string
	// The boot environment that the machine should boot into.  This
	// must be the name of a boot environment present in the backend.
	// If this field is not present or blank, the global default bootenv
	// will be used instead.
	BootEnv string
	// An array of profiles to apply to this machine in order when looking
	// for a parameter during rendering.
	Profiles []string
	// The Parameters that have been directly set on the Machine.
	Params map[string]interface{}
	// The tasks this machine has to run.
	Tasks []string
	// The index into the Tasks list for the task that is currently
	// running (if a task is running) or the next task that will run (if
	// no task is currently running).  If -1, then the first task will
	// run next, and if it is equal to the length of the Tasks list then
	// all the tasks have finished running.
	//
	// required: true
	CurrentTask int
	// This tracks the number of retry attempts for the current task.
	// When a task succeeds, the retry value is reset.
	RetryTaskAttempt int
	// This list of previous task lists and current tasks to handle errors.
	// Upon completing the list, the previous task list will be executed.
	//
	// This will be capped to a depth of 1.  Error failures are not handled can not be handled.
	TaskErrorStacks []*TaskStack
	// Indicates if the machine can run jobs or not.  Failed jobs mark the machine
	// not runnable.
	//
	// required: true
	Runnable bool
	// Secret for machine token revocation.  Changing the secret will invalidate
	// all existing tokens for this machine
	Secret string `index:",ignore"`
	// OS is the operating system that the node is running in.  It is updated by Sledgehammer and by
	// the various OS install tasks.
	//
	OS string
	// HardwareAddrs is a list of MAC addresses we expect that the system might boot from.
	// This must be filled out to enable MAC address based booting from the various bootenvs,
	// and must be updated if the MAC addresses for a system change for whatever reason.
	//
	HardwareAddrs []string
	// Workflow is the workflow that is currently responsible for processing machine tasks.
	//
	// required: true
	Workflow string
	// Arch is the machine architecture. It should be an arch that can
	// be fed into $GOARCH.
	//
	// required: true
	Arch string
	// Locked indicates that changes to the Machine by users are not
	// allowed, except for unlocking the machine, which will always
	// generate an Audit event.
	//
	// required: true
	Locked bool
	// Contexts contains the name of the current execution context for the machine.
	// An empty string indicates that the agent running on the machine should be executing tasks,
	// and any other value means that an agent running with its context set for this value should
	// be executing tasks.
	Context string
	// Fingerprint is a collection of data that can (in theory) be used to uniquely identify
	// a machine based on various DMI information.  This (in conjunction with HardwareAddrs)
	// is used to uniquely identify a Machine using a score based on how many total items in the Fingerprint
	// match.
	Fingerprint MachineFingerprint `index:",ignore"`
	// Pool contains the pool the machine is in.
	// Unset machines will join the default Pool
	Pool string
	// PoolAllocated defines if the machine is allocated in this pool
	// This is a calculated field.
	PoolAllocated bool
	// PoolStatus contains the status of this machine in the Pool.
	//    Values are defined in Pool.PoolStatuses
	PoolStatus PoolStatus
	// WorkflowCopmlete indicates if the workflow is complete
	WorkflowComplete bool
}

func (n *Machine) IsLocked() bool {
	return n.Locked
}

func (n *Machine) GetMeta() Meta {
	return n.Meta
}

func (n *Machine) SetMeta(d Meta) {
	n.Meta = d
}

func (n *Machine) Validate() {
	if arch, ok := SupportedArch(n.Arch); !ok {
		n.Errorf("Unsupported arch %s", n.Arch)
	} else if arch != n.Arch {
		n.Errorf("Please use %s for Arch instead of %s", arch, n.Arch)
	}
	n.AddError(ValidMachineName("Invalid Name", n.Name))
	n.AddError(ValidName("Invalid Stage", n.Stage))
	n.AddError(ValidName("Invalid BootEnv", n.BootEnv))
	if n.Workflow != "" {
		n.AddError(ValidName("Invalid Workflow", n.Workflow))
	}
	for _, p := range n.Profiles {
		n.AddError(ValidName("Invalid Profile", p))
	}
	for _, t := range n.Tasks {
		parts := strings.SplitN(t, ":", 2)
		if len(parts) == 2 {
			switch parts[0] {
			case "stage":
				n.AddError(ValidName("Invalid Stage", parts[1]))
			case "bootenv":
				n.AddError(ValidName("Invalid BootEnv", parts[1]))
			case "chroot", "context":
			case "action":
				pparts := strings.SplitN(parts[1], ":", 2)
				if len(pparts) == 2 {
					n.AddError(ValidName("Invalid Plugin", pparts[0]))
					n.AddError(ValidName("Invalid Action", pparts[1]))
				} else {
					n.AddError(ValidName("Invalid Action", parts[1]))
				}
			}
		} else {
			n.AddError(ValidName("Invalid Task", t))
		}
	}
	for _, m := range n.HardwareAddrs {
		if _, err := net.ParseMAC(m); err != nil {
			n.Errorf("Invalid Hardware Address `%s`: %v", m, err)
		}
	}
}

func (n *Machine) UUID() string {
	return n.Uuid.String()
}

func (n *Machine) Prefix() string {
	return "machines"
}

func (n *Machine) Key() string {
	return n.UUID()
}

func (n *Machine) KeyName() string {
	return "Uuid"
}

func (m *Machine) GetDescription() string {
	return m.Description
}

func (n *Machine) Fill() {
	if n.Meta == nil {
		n.Meta = Meta{}
	}
	n.Validation.fill(n)
	if n.Profiles == nil {
		n.Profiles = []string{}
	}
	if n.Tasks == nil {
		n.Tasks = []string{}
	}
	if n.Params == nil {
		n.Params = map[string]interface{}{}
	}
	if n.HardwareAddrs == nil {
		n.HardwareAddrs = []string{}
	}
	if n.TaskErrorStacks == nil {
		n.TaskErrorStacks = []*TaskStack{}
	}
	if n.Arch == "" {
		n.Arch = "amd64"
	}
	(&n.Fingerprint).Fill()
}

func (n *Machine) AuthKey() string {
	return n.Key()
}

func (b *Machine) SliceOf() interface{} {
	s := []*Machine{}
	return &s
}

func (b *Machine) ToModels(obj interface{}) []Model {
	items := obj.(*[]*Machine)
	res := make([]Model, len(*items))
	for i, item := range *items {
		res[i] = Model(item)
	}
	return res
}

// match Param interface
func (b *Machine) GetParams() map[string]interface{} {
	return copyMap(b.Params)
}

func (b *Machine) SetParams(p map[string]interface{}) {
	b.Params = copyMap(p)
}

// match Profiler interface
func (b *Machine) GetProfiles() []string {
	return b.Profiles
}

func (b *Machine) SetProfiles(p []string) {
	b.Profiles = p
}

// match BootEnver interface
func (b *Machine) GetBootEnv() string {
	return b.BootEnv
}

func (b *Machine) SetBootEnv(s string) {
	b.BootEnv = s
}

// match TaskRunner interface
func (b *Machine) GetTasks() []string {
	return b.Tasks
}

func (b *Machine) SetTasks(t []string) {
	b.Tasks = t
}

func (b *Machine) RunningTask() int {
	return b.CurrentTask
}

func (b *Machine) SetName(n string) {
	b.Name = n
}

// SplitTasks slits the machine's Tasks list into 3 subsets:
//
// 1. the immutable past, which cannot be chnaged by task list modification
//
// 2. The mutable present, which contains tasks that can be deleted, and where tasks can be added.
//
// 3. The immutable future, which also cannot be changed.
func (b *Machine) SplitTasks() (thePast []string, thePresent []string, theFuture []string) {
	thePast, thePresent, theFuture = []string{}, []string{}, []string{}
	if len(b.Tasks) == 0 {
		return
	}
	if b.CurrentTask == -1 {
		thePresent = b.Tasks[:]
	} else if b.CurrentTask >= len(b.Tasks) {
		thePast = b.Tasks[:]
	} else {
		thePast = b.Tasks[:b.CurrentTask+1]
		thePresent = b.Tasks[b.CurrentTask+1:]
	}
	for i := 0; i < len(thePresent); i++ {
		if strings.HasPrefix(thePresent[i], "stage:") {
			theFuture = thePresent[i:]
			thePresent = thePresent[:i]
			break
		}
	}
	return
}

// AddTasks is a helper for adding tasks to the machine Tasks list in
// the mutable present.
func (b *Machine) AddTasks(offset int, tasks ...string) error {
	thePast, thePresent, theFuture := b.SplitTasks()
	if offset < 0 {
		offset += len(thePresent) + 1
		if offset < 0 {
			offset = len(thePresent)
		}
	}
	if offset >= len(thePresent) {
		offset = len(thePresent)
	}
	if offset == 0 {
		if len(thePresent) >= (len(tasks)+offset) &&
			reflect.DeepEqual(tasks, thePresent[offset:offset+len(tasks)]) {
			// We are already in the desired task state.
			return nil
		}
		thePresent = append(tasks, thePresent...)
	} else if offset == len(thePresent) {
		if len(thePresent) >= len(tasks) &&
			reflect.DeepEqual(tasks, thePresent[len(thePresent)-len(tasks):]) {
			// We are alredy in the desired state
			return nil
		}
		thePresent = append(thePresent, tasks...)
	} else {
		if len(thePresent[offset:]) >= len(tasks) &&
			reflect.DeepEqual(tasks, thePresent[offset:offset+len(tasks)]) {
			// Already in the desired state
			return nil
		}
		res := []string{}
		res = append(res, thePresent[:offset]...)
		res = append(res, tasks...)
		res = append(res, thePresent[offset:]...)
		thePresent = res
	}
	thePresent = append(thePresent, theFuture...)
	b.Tasks = append(thePast, thePresent...)
	return nil
}

// DelTasks allows you to delete tasks in the mutable present.
func (b *Machine) DelTasks(tasks ...string) {
	if len(tasks) == 0 {
		return
	}
	thePast, thePresent, theFuture := b.SplitTasks()
	if len(thePresent) == 0 {
		return
	}
	nextThePresent := []string{}
	i := 0
	for _, c := range thePresent {
		if i < len(tasks) && tasks[i] == c {
			i++
		} else {
			nextThePresent = append(nextThePresent, c)
		}
	}
	nextThePresent = append(nextThePresent, theFuture...)
	b.Tasks = append(thePast, nextThePresent...)
}

func (b *Machine) CanHaveActions() bool {
	return true
}
