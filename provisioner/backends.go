package provisioner

import (
	"encoding/json"
	"fmt"

	"github.com/digitalrebar/digitalrebar/go/common/store"
)

type keySaver interface {
	prefix() string
	key() string
	typeName() string
	onChange(interface{}) error
	onDelete() error
	newIsh() keySaver
	RebuildRebarData() error
}

func registerBackends(s store.SimpleStore) {
	backendMux.Lock()
	defer backendMux.Unlock()
	t := &Template{}
	b := &BootEnv{}
	m := &Machine{}
	backends[t.prefix()] = s.Sub(t.prefix())
	backends[b.prefix()] = s.Sub(b.prefix())
	backends[m.prefix()] = s.Sub(m.prefix())
}

func getBackend(t keySaver) store.SimpleStore {
	backendMux.Lock()
	defer backendMux.Unlock()
	res, ok := backends[t.prefix()]
	if !ok {
		Logger.Fatalf("%s: No registered storage backend!", t.prefix())
	}
	return res
}

func list(t keySaver) [][]byte {
	backend := getBackend(t)

	keys, err := backend.Keys()
	if err != nil {
		Logger.Fatalf("%s: Error getting keys: %v", t.prefix(), err)
	}
	res := make([][]byte, len(keys))
	for i, k := range keys {
		res[i], err = backend.Load(k)
		if err != nil {
			Logger.Fatalf("%s: Error reading contents for %s: %v", t.prefix(), k, err)
		}
	}
	return res
}

//	save(keySaver, interface{}) error
//      remove(keySaver) error

func load(t keySaver) error {
	backend := getBackend(t)
	buf, err := backend.Load(t.key())
	if err != nil {
		return fmt.Errorf("%s: Failed to load %s: %v", t.prefix(), t.key(), err)
	}
	return json.Unmarshal(buf, &t)
}

func save(newThing keySaver, oldThing interface{}) error {
	backend := getBackend(newThing)

	if err := newThing.onChange(oldThing); err != nil && oldThing != nil {
		return err
	}
	buf, err := json.Marshal(newThing)
	if err != nil {
		return fmt.Errorf("%s: Failed to marshal %s: %v", newThing.prefix(), newThing.key(), err)
	}
	return backend.Save(newThing.key(), buf)
}

func remove(t keySaver) error {
	backend := getBackend(t)
	if err := load(t); err != nil {
		return err
	}
	if err := t.onDelete(); err != nil {
		return err
	}
	return backend.Remove(t.key())
}
