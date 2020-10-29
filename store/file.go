package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

type File struct {
	storeBase
	Path string
	data struct {
		Sections map[string]map[string]interface{} `json:"sections,omitempty"`
		Meta     map[string]string                 `json:"meta,omitempty"`
	}
}

func (f *File) Type() string {
	return "file"
}

func (f *File) MetaData() map[string]string {
	f.RLock()
	defer f.RUnlock()
	res := map[string]string{}
	for k, v := range f.data.Meta {
		res[k] = v
	}
	return res
}

func (f *File) SetMetaData(vals map[string]string) error {
	f.Lock()
	defer f.Unlock()
	oldMeta := f.data.Meta
	f.data.Meta = map[string]string{}
	for k, v := range vals {
		f.data.Meta[k] = v
	}
	err := f.save()
	if err != nil {
		f.data.Meta = oldMeta
	}
	if n, ok := vals["Name"]; ok {
		f.name = n
	}
	return err
}

func (f *File) Open(codec Codec) error {
	if f.Path == "" {
		return fmt.Errorf("Cannot store data at ''")
	}
	fullPath, err := filepath.Abs(filepath.Clean(f.Path))
	if err != nil {
		return err
	}
	if codec == nil {
		codec = DefaultCodec
	}
	f.Codec = codec
	if err := os.MkdirAll(path.Dir(fullPath), 0755); err != nil {
		return err
	}

	buf, err := ioutil.ReadFile(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	f.data.Meta = map[string]string{}
	f.data.Sections = map[string]map[string]interface{}{}
	if buf != nil {
		if err := f.Decode(buf, &f.data); err != nil {
			return err
		}
	}
	f.name = f.data.Meta["Name"]
	f.opened = true
	return nil
}

func (f *File) Prefixes() ([]string, error) {
	f.RLock()
	defer f.RUnlock()
	res := []string{}
	for k := range f.data.Sections {
		res = append(res, k)
	}
	return res, nil
}

func (f *File) Keys(prefix string) ([]string, error) {
	f.RLock()
	defer f.RUnlock()
	f.panicIfClosed()
	vals, ok := f.data.Sections[prefix]
	if !ok {
		return []string{}, nil
	}
	res := make([]string, 0, len(vals))
	for k := range vals {
		res = append(res, k)
	}
	return res, nil
}

func (f *File) Exists(prefix, key string) bool {
	f.RLock()
	defer f.RUnlock()
	f.panicIfClosed()
	vals, ok := f.data.Sections[prefix]
	if !ok {
		return ok
	}
	_, ok = vals[key]
	return ok
}

func (f *File) Load(prefix, key string, val interface{}) error {
	f.RLock()
	defer f.RUnlock()
	f.panicIfClosed()
	if !f.Exists(prefix, key) {
		return os.ErrNotExist
	}
	if err := remarshal(f.data.Sections[prefix][key], &val); err != nil {
		return err
	}
	if ro, ok := val.(ReadOnlySetter); ok {
		ro.SetReadOnly(f.ReadOnly())
	}
	if bb, ok := val.(BundleSetter); ok {
		n := f.Name()
		if n != "" {
			bb.SetBundle(n)
		}
	}
	return nil
}

func (f *File) save() error {
	f.panicIfClosed()
	buf, err := f.Encode(f.data)
	if err != nil {
		return err
	}
	return safeReplace(f.Path, buf)
}

func (f *File) Save(prefix, key string, val interface{}) error {
	f.Lock()
	defer f.Unlock()
	if f.readOnly {
		return UnWritable(key)
	}
	if _, ok := f.data.Sections[prefix]; !ok {
		f.data.Sections[prefix] = map[string]interface{}{}
	}
	f.data.Sections[prefix][key] = val
	return f.save()
}

func (f *File) Remove(prefix, key string) error {
	f.Lock()
	defer f.Unlock()
	if f.readOnly {
		return UnWritable(key)
	}
	if _, ok := f.data.Sections[prefix]; !ok {
		return os.ErrNotExist
	}
	if _, ok := f.data.Sections[prefix][key]; !ok {
		return os.ErrNotExist
	}
	delete(f.data.Sections[prefix], key)
	return f.save()
}
