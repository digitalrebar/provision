package provisioner

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	consul "github.com/hashicorp/consul/api"
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

type storageBackend interface {
	list(keySaver) [][]byte
	save(keySaver, interface{}) error
	load(keySaver) error
	remove(keySaver) error
}

type fileBackend string

func newFileBackend(path string) (fileBackend, error) {
	fullPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return "", err
	}
	return fileBackend(fullPath), nil
}

func (f fileBackend) mkThingPath(thing keySaver) string {
	fullPath := filepath.Join(string(f), thing.prefix())
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		Logger.Fatalf("file: Cannot create %s: %v", fullPath, err)
	}
	return fullPath
}

func (f fileBackend) mkThingName(thing keySaver) string {
	return filepath.Join(string(f), thing.key()) + ".json"
}

func (f fileBackend) list(thing keySaver) [][]byte {
	dir := f.mkThingPath(thing)
	file, err := os.Open(dir)
	if err != nil {
		Logger.Fatalf("file: Failed to open dir %s: %v", dir, err)
	}
	names, err := file.Readdirnames(0)
	if err != nil {
		Logger.Fatalf("file: Failed to get listing for dir %s: %v", dir, err)
	}
	res := make([][]byte, 0, len(names))
	for _, name := range names {
		if !strings.HasSuffix(name, ".json") {
			continue
		}
		fullName := filepath.Join(dir, name)
		buf, err := ioutil.ReadFile(fullName)
		if err != nil {
			Logger.Fatalf("file: Failed to read info for %s: %v", fullName, err)
		}
		res = append(res, buf)
	}
	return res
}

func (f fileBackend) load(thing keySaver) error {
	fullName := f.mkThingName(thing)
	buf, err := ioutil.ReadFile(fullName)
	if err != nil {
		return fmt.Errorf("file: Failed to read %s: %v", fullName, err)
	}
	return json.Unmarshal(buf, &thing)
}

func (f fileBackend) save(newThing keySaver, oldThing interface{}) error {
	f.mkThingPath(newThing)
	if err := newThing.onChange(oldThing); err != nil && oldThing != nil {
		return err
	}
	fullPath := f.mkThingName(newThing)
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("file: Failed to open thing %s: %v", fullPath, err)
	}
	enc := json.NewEncoder(file)
	if err := enc.Encode(newThing); err != nil {
		os.Remove(fullPath)
		file.Close()
		return fmt.Errorf("file: Failed to save %s: %v", fullPath, err)
	}
	file.Sync()
	file.Close()
	return nil
}

func (f fileBackend) remove(thing keySaver) error {
	if err := f.load(thing); err != nil {
		return err
	}
	if err := thing.onDelete(); err != nil {
		return err
	}
	return os.Remove(f.mkThingName(thing))
}

type consulBackend struct {
	kv      *consul.KV
	baseKey string
}

func (cb *consulBackend) makePrefix(thing keySaver) string {
	return path.Join(cb.baseKey, thing.prefix())
}

func (cb *consulBackend) makeKey(thing keySaver) string {
	return path.Join(cb.baseKey, thing.key())
}

func newConsulBackend(baseKey string) (*consulBackend, error) {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}
	backend := &consulBackend{
		kv:      client.KV(),
		baseKey: baseKey,
	}
	return backend, nil
}

func (cb *consulBackend) list(thing keySaver) [][]byte {
	keypairs, _, err := cb.kv.List(cb.makePrefix(thing), nil)
	if err != nil {
		return [][]byte{}
	}
	res := make([][]byte, len(keypairs))
	for i, kp := range keypairs {
		res[i] = kp.Value
	}
	return res
}

func (cb *consulBackend) save(newThing keySaver, oldThing interface{}) error {
	if err := newThing.onChange(oldThing); err != nil && oldThing != nil {
		return err
	}
	buf, err := json.Marshal(newThing)
	if err != nil {
		return fmt.Errorf("consul: Failed to marshal %+v: %v", newThing, err)
	}
	kp := &consul.KVPair{Value: buf, Key: cb.makeKey(newThing)}
	if _, err := cb.kv.Put(kp, nil); err != nil {
		return fmt.Errorf("consul: Failed to save %s: %v", kp.Key, err)
	}
	err = newThing.RebuildRebarData()
	return err
}

func (cb *consulBackend) load(s keySaver) error {
	key := cb.makeKey(s)
	kp, _, err := cb.kv.Get(key, nil)
	if err != nil {
		return fmt.Errorf("consul: Communication failure: %v", err)
	} else if kp == nil {
		return fmt.Errorf("consul: Failed to load %v", key)
	}
	if err := json.Unmarshal(kp.Value, &s); err != nil {
		return fmt.Errorf("consul: Failed to unmarshal %s: %v", kp.Key, err)
	}
	return nil
}

func (cb *consulBackend) remove(s keySaver) error {
	if err := cb.load(s); err != nil {
		return err
	}
	if err := s.onDelete(); err != nil {
		return err
	}
	key := cb.makeKey(s)
	if _, err := cb.kv.Delete(key, nil); err != nil {
		return fmt.Errorf("consul: Failed to delete %v: %v", key, err)
	}
	err := s.RebuildRebarData()
	return err
}
