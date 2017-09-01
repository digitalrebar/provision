package models

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

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

func New(kind string) (Model, error) {
	switch kind {
	case "stages", "stage":
		return &Stage{}, nil
	case "bootenvs", "bootenv":
		return &BootEnv{}, nil
	case "jobs", "job":
		return &Job{}, nil
	case "leases", "lease":
		return &Lease{}, nil
	case "machines", "machine":
		return &Machine{}, nil
	case "params", "param":
		return &Param{}, nil
	case "plugins", "plugin":
		return &Plugin{}, nil
	case "prefs", "pref":
		return &Pref{}, nil
	case "profiles", "profile":
		return &Profile{}, nil
	case "reservations", "reservation":
		return &Reservation{}, nil
	case "subnets", "subnet":
		return &Subnet{}, nil
	case "tasks", "task":
		return &Task{}, nil
	case "templates", "template":
		return &Template{}, nil
	case "users", "user":
		return &User{}, nil
	default:
		return nil, fmt.Errorf("No such Model: %s", kind)
	}
}

func Clone(m Model) (Model, error) {
	res, err := New(m.Prefix())
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	enc, dec := gob.NewEncoder(&buf), gob.NewDecoder(&buf)
	err = enc.Encode(m)
	if err != nil {
		return nil, err
	}
	err = dec.Decode(res)
	return res, err
}
