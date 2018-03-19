package models

import (
	"net"
	"reflect"
	"strings"

	"github.com/pborman/uuid"
)

// Machine represents a single bare-metal system that the provisioner
// should manage the boot environment for.
// swagger:model
type Machine struct {
	Validation
	Access
	Meta
	// The name of the machine.  THis must be unique across all
	// machines, and by convention it is the FQDN of the machine,
	// although nothing enforces that.
	//
	// required: true
	// swagger:strfmt hostname
	Name string
	// A description of this machine.  This can contain any reference
	// information for humans you want associated with the machine.
	Description string
	// The UUID of the machine.
	// This is auto-created at Create time, and cannot change afterwards.
	//
	// required: true
	// swagger:strfmt uuid
	Uuid uuid.UUID
	// The UUID of the job that is currently running on the machine.
	//
	// swagger:strfmt uuid
	CurrentJob uuid.UUID
	// The IPv4 address of the machine that should be used for PXE
	// purposes.  Note that this field does not directly tie into DHCP
	// leases or reservations -- the provisioner relies solely on this
	// address when determining what to render for a specific machine.
	//
	// swagger:strfmt ipv4
	Address net.IP
	// An optional value to indicate tasks and profiles to apply.
	Stage string
	// The boot environment that the machine should boot into.  This
	// must be the name of a boot environment present in the backend.
	// If this field is not present or blank, the global default bootenv
	// will be used instead.
	BootEnv string
	// An array of profiles to apply to this machine in order when looking
	// for a parameter during rendering.
	Profiles []string
	//
	// The Machine specific Profile Data - only used for the map (name and other
	// fields not used - THIS IS DEPRECATED AND WILL GO AWAY.
	// Data will migrated from this struct to Params and then cleared.
	Profile Profile
	// Replaces the Profile.
	Params map[string]interface{}
	// The tasks this machine has to run.
	Tasks []string
	// required: true
	CurrentTask int
	// Indicates if the machine can run jobs or not.  Failed jobs mark the machine
	// not runnable.
	//
	// required: true
	Runnable bool

	// Secret for machine token revocation.  Changing the secret will invalidate
	// all existing tokens for this machine
	Secret string
	// OS is the operating system that the node is running in
	//
	OS string
	// HardwareAddrs is a list of MAC addresses we expect that the system might boot from.
	//
	//
	HardwareAddrs []string
	// Workflow is the workflow that is currently responsible for processing machine tasks.
	//
	// required: true
	Workflow string
}

func (n *Machine) Validate() {
	n.AddError(ValidName("Invalid Name", n.Name))
	n.AddError(ValidName("Invalid Stage", n.Stage))
	n.AddError(ValidName("Invalid BootEnv", n.BootEnv))
	for _, p := range n.Profiles {
		n.AddError(ValidName("Invalid Profile", p))
	}
	for _, t := range n.Tasks {
		if n.Workflow == "" {
			n.AddError(ValidName("Invalid Task", t))
		} else {
			parts := strings.SplitN(t, ":", 2)
			if len(parts) == 2 {
				if parts[0] != "stage" {
					n.Errorf("Invalid Task Step %s", t)
				} else {
					n.AddError(ValidName("Invalid Stage", parts[1]))
				}
			} else {
				n.AddError(ValidName("Invalid Task", t))
			}
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

func (n *Machine) Fill() {
	if n.Meta == nil {
		n.Meta = Meta{}
	}
	n.Validation.fill()
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

// match Paramer interface
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

func (b *Machine) SplitTasks() (thePast []string, thePresent []string, theFuture []string) {
	thePast, thePresent, theFuture = []string{}, []string{}, []string{}
	if len(b.Tasks) == 0 {
		return
	}
	if b.CurrentTask == -1 {
		thePresent = b.Tasks[:]
	} else if b.CurrentTask == len(b.Tasks) {
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
