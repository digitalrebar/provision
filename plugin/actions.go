package plugin

import (
	"fmt"
	"net"
	"sort"
	"sync"
	"time"

	"github.com/pborman/uuid"
)

// Plugins can provide actions for machines
// Assumes that there are parameters on the
// call in addition to the machine.
//
// swagger:model
type AvailableAction struct {
	Provider       string
	Command        string
	RequiredParams []string
	OptionalParams []string

	plugin *RunningPlugin
	ma     *MachineActions

	lock      sync.Mutex
	inflight  int
	unloading bool
}

//
// Params is built from the caller, plus
// the machine, plus profiles, plus global.
//
type MachineAction struct {
	Name    string
	Uuid    uuid.UUID
	Address net.IP
	BootEnv string
	Command string
	Params  map[string]interface{}
}

type MachineActions struct {
	actions map[string]*AvailableAction
	lock    sync.Mutex
}

func NewMachineActions() *MachineActions {
	return &MachineActions{actions: make(map[string]*AvailableAction, 0)}
}

func (ma *MachineActions) Add(aa *AvailableAction) error {
	ma.lock.Lock()
	defer ma.lock.Unlock()

	if _, ok := ma.actions[aa.Command]; ok {
		return fmt.Errorf("Duplicate Action %s: already present\n", aa.Command)
	}
	ma.actions[aa.Command] = aa
	aa.ma = ma
	return nil
}

func (ma *MachineActions) Remove(aa *AvailableAction) error {
	var err error
	ma.lock.Lock()
	if _, ok := ma.actions[aa.Command]; !ok {
		err = fmt.Errorf("Missing Action %s: already removed\n", aa.Command)
	} else {
		delete(ma.actions, aa.Command)
	}
	ma.lock.Unlock()

	aa.Unload()
	return err
}

func (ma *MachineActions) List() []*AvailableAction {
	ma.lock.Lock()
	defer ma.lock.Unlock()

	// get the list of keys and sort them
	keys := []string{}
	for key := range ma.actions {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	answer := []*AvailableAction{}
	for _, key := range keys {
		answer = append(answer, ma.actions[key])
	}
	return answer

}

func (ma *MachineActions) Get(name string) (a *AvailableAction, ok bool) {
	ma.lock.Lock()
	defer ma.lock.Unlock()

	a, ok = ma.actions[name]
	return
}

func (aa *AvailableAction) Reserve() error {
	aa.lock.Lock()
	defer aa.lock.Unlock()

	if aa.unloading {
		return fmt.Errorf("Action not available %s: unloading\n", aa.Command)
	}
	aa.inflight += 1
	return nil
}

func (aa *AvailableAction) Release() {
	aa.lock.Lock()
	defer aa.lock.Unlock()

	aa.inflight -= 1
}

func (aa *AvailableAction) Unload() {
	aa.lock.Lock()
	aa.unloading = true
	for aa.inflight != 0 {
		aa.lock.Unlock()
		time.Sleep(time.Millisecond * 15)
		aa.lock.Lock()
	}
	aa.lock.Unlock()
	return
}

func (aa *AvailableAction) Run(maa *MachineAction) error {
	if err := aa.Reserve(); err != nil {
		return nil
	}
	defer aa.Release()

	return aa.plugin.Client.Action(maa)
}
