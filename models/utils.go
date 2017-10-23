package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// swagger:model
type BlobInfo struct {
	Path string
	Size int64
}

type Model interface {
	Prefix() string
	Key() string
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
		&Job{},
		&Lease{},
		&Machine{},
		&Param{},
		&Plugin{},
		&Pref{},
		&Profile{},
		&Reservation{},
		&Subnet{},
		&Task{},
		&Template{},
		&User{},
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
