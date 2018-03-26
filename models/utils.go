package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"
)

func copyMap(m map[string]interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	for k, v := range m {
		res[k] = v
	}
	return res
}

// swagger:model
type BlobInfo struct {
	Path string
	Size int64
}

type Model interface {
	Prefix() string
	Key() string
	KeyName() string
}

type Filler interface {
	Model
	Fill()
}

type Slicer interface {
	Filler
	SliceOf() interface{}
	ToModels(interface{}) []Model
}

func All() []Model {
	return []Model{
		&BootEnv{},
		&Interface{},
		&Job{},
		&Lease{},
		&Machine{},
		&Param{},
		&Plugin{},
		&PluginProvider{},
		&Pref{},
		&Profile{},
		&Reservation{},
		&Stage{},
		&Subnet{},
		&Task{},
		&Template{},
		&User{},
		&Workflow{},
	}
}

func New(kind string) (Slicer, error) {
	var res Slicer
	switch kind {
	case "bootenvs", "bootenv":
		res = &BootEnv{}
	case "interfaces":
		res = &Interface{}
	case "jobs", "job":
		res = &Job{}
	case "leases", "lease":
		res = &Lease{}
	case "machines", "machine":
		res = &Machine{}
	case "params", "param":
		res = &Param{}
	case "plugins", "plugin":
		res = &Plugin{}
	case "plugin_providers", "plugin_provider":
		res = &PluginProvider{}
	case "preferences", "preference":
		res = &Pref{}
	case "profiles", "profile":
		res = &Profile{}
	case "reservations", "reservation":
		res = &Reservation{}
	case "stages", "stage":
		res = &Stage{}
	case "subnets", "subnet":
		res = &Subnet{}
	case "tasks", "task":
		res = &Task{}
	case "templates", "template":
		res = &Template{}
	case "users", "user":
		res = &User{}
	case "workflows", "workflow":
		res = &Workflow{}
	default:
		return nil, fmt.Errorf("No such Model: %s", kind)
	}
	res.Fill()
	return res, nil
}

func Clone(m Model) Model {
	if m == nil {
		return nil
	}
	res, err := New(m.Prefix())
	if err != nil {
		log.Panicf("Failed to make a new %s: %v", m.Prefix(), err)
	}
	buf := bytes.Buffer{}
	enc, dec := json.NewEncoder(&buf), json.NewDecoder(&buf)
	if err := enc.Encode(m); err != nil {
		log.Panicf("Failed to encode %s:%s: %v", m.Prefix(), m.Key(), err)
	}
	if err := dec.Decode(res); err != nil {
		log.Panicf("Failed to decode %s:%s: %v", m.Prefix(), m.Key(), err)
	}
	return res
}

var (
	validName      = regexp.MustCompile(`^\pL+([- _.]+|\pN+|\pL+)+$`)
	validParamName = regexp.MustCompile(`^\pL+([- _./]+|\pN+|\pL+)+$`)
)

func ValidName(msg, s string) error {
	if validName.MatchString(s) {
		return nil
	}
	return fmt.Errorf("%s `%s`", msg, s)
}

func ValidParamName(msg, s string) error {
	if validParamName.MatchString(s) {
		return nil
	}
	return fmt.Errorf("%s `%s`", msg, s)
}

type NameSetter interface {
	Model
	SetName(string)
}

type Paramer interface {
	Model
	GetParams() map[string]interface{}
	SetParams(map[string]interface{})
}

type Profiler interface {
	Model
	GetProfiles() []string
	SetProfiles([]string)
}

type BootEnver interface {
	Model
	GetBootEnv() string
	SetBootEnv(string)
}

type Tasker interface {
	Model
	GetTasks() []string
	SetTasks([]string)
}

type TaskRunner interface {
	Tasker
	RunningTask() int
}

// Only implement this if you want actions
type Actor interface {
	Model
	CanHaveActions() bool
}

func FibBackoff(thunk func() error) {
	timeouts := []time.Duration{
		time.Second,
		time.Second,
		2 * time.Second,
		3 * time.Second,
		5 * time.Second,
		8 * time.Second,
	}
	for _, d := range timeouts {
		if thunk() == nil {
			return
		}
		time.Sleep(d)
	}
}
