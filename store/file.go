package store

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
)

type File struct {
	storeBase
	Path string
	vals map[string][]byte
	meta map[string]string
}

func (f *File) Type() string {
	return "file"
}

func (f *File) MetaData() map[string]string {
	f.RLock()
	defer f.RUnlock()
	if f.parentStore != nil {
		return f.parentStore.(*File).MetaData()
	}
	res := map[string]string{}
	for k, v := range f.meta {
		res[k] = v
	}
	return res
}

func (f *File) SetMetaData(vals map[string]string) error {
	f.Lock()
	defer f.Unlock()
	if f.parentStore != nil {
		return f.parentStore.(*File).SetMetaData(vals)
	}
	oldMeta := f.meta
	f.meta = map[string]string{}
	for k, v := range vals {
		f.meta[k] = v
	}
	err := f.save()
	if err != nil {
		f.meta = oldMeta
	}
	if n, ok := vals["Name"]; ok {
		f.name = n
	}
	return err
}

func (f *File) MakeSub(path string) (Store, error) {
	f.Lock()
	defer f.Unlock()
	f.panicIfClosed()
	if child, ok := f.subStores[path]; ok {
		return child, nil
	}
	sub := &File{}
	sub.Codec = f.Codec
	sub.vals = map[string][]byte{}
	sub.opened = true
	addSub(f, sub, path)
	return sub, nil
}

func (f *File) mux() *sync.RWMutex {
	f.RLock()
	defer f.RUnlock()
	if f.parentStore != nil {
		return f.parentStore.(*File).mux()
	}
	return &f.RWMutex
}

func (f *File) open(vals map[string]interface{}) error {
	f.vals = map[string][]byte{}
	f.meta = map[string]string{}
	for k, v := range vals {
		switch k {
		case "sections":
			subSections, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Invalid sections declaration: %#v", v)
			}
			for subName, subVals := range subSections {
				sub := &File{}
				sub.Codec = f.Codec
				if err := sub.open(subVals.(map[string]interface{})); err != nil {
					return err
				}
				addSub(f, sub, subName)
			}
		case "meta":
			metaData, ok := v.(map[string]interface{})
			if !ok {
				return fmt.Errorf("Invalid metadata declaration: %#v", v)
			}
			for metaName, metaVal := range metaData {
				if val, ok := metaVal.(string); ok {
					f.meta[metaName] = val
				} else {
					return fmt.Errorf("Metadata value %#v is not a string", metaVal)
				}
			}
		default:
			buf, err := f.Encode(v)
			if err != nil {
				return err
			}
			f.vals[k] = buf
		}
	}
	f.opened = true
	md := f.MetaData()
	if n, ok := md["Name"]; ok {
		f.name = n
	}
	return nil
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
	vals := map[string]interface{}{}
	if err := os.MkdirAll(path.Dir(fullPath), 0755); err != nil {
		return err
	}

	buf, err := ioutil.ReadFile(fullPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if buf != nil {
		if err := f.Decode(buf, &vals); err != nil {
			return err
		}
	}
	return f.open(vals)
}

func (f *File) Keys() ([]string, error) {
	mux := f.mux()
	mux.RLock()
	defer mux.RUnlock()
	f.panicIfClosed()
	res := make([]string, 0, len(f.vals))
	for k := range f.vals {
		res = append(res, k)
	}
	return res, nil
}

func (f *File) Load(key string, val interface{}) error {
	mux := f.mux()
	mux.RLock()
	defer mux.RUnlock()
	f.panicIfClosed()
	buf, ok := f.vals[key]
	if ok {
		if err := f.Decode(buf, val); err != nil {
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
	return os.ErrNotExist
}

func (f *File) prepSave() (map[string]interface{}, error) {
	res := map[string]interface{}{}
	for k, v := range f.vals {
		var obj interface{}
		if err := f.Decode(v, &obj); err != nil {
			return nil, err
		}
		res[k] = obj
	}
	if len(f.subStores) > 0 {
		subs := map[string]interface{}{}
		for subName, subStore := range f.subStores {
			subVals, err := subStore.(*File).prepSave()
			if err != nil {
				return nil, err
			}
			subs[subName] = subVals
		}
		res["sections"] = subs
	}
	if len(f.meta) > 0 {
		res["meta"] = f.meta
	}
	return res, nil
}

func (f *File) save() error {
	f.panicIfClosed()
	if f.parentStore != nil {
		parent := f.parentStore.(*File)
		return parent.save()
	}
	toSave, err := f.prepSave()
	if err != nil {
		return err
	}
	buf, err := f.Encode(toSave)
	if err != nil {
		return err
	}
	return safeReplace(f.Path, buf)
}

func (f *File) Save(key string, val interface{}) error {
	mux := f.mux()
	mux.Lock()
	defer mux.Unlock()
	if f.readOnly {
		return UnWritable(key)
	}
	buf, err := f.Encode(val)
	if err != nil {
		return err
	}
	f.vals[key] = buf
	return f.save()
}

func (f *File) Remove(key string) error {
	mux := f.mux()
	mux.Lock()
	defer mux.Unlock()
	if f.readOnly {
		return UnWritable(key)
	}
	if _, ok := f.vals[key]; !ok {
		return os.ErrNotExist
	}
	delete(f.vals, key)
	return f.save()
}
