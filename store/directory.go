package store

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Directory implements a Store that is backed by a local directory tree.
type Directory struct {
	storeBase
	Path string
}

func (d *Directory) Type() string {
	return "directory"
}

func (f *Directory) filename(n string) string {
	return filepath.Join(f.Path, url.QueryEscape(n))
}

func (d *Directory) MetaData() (res map[string]string) {
	d.RLock()
	defer d.RUnlock()
	if d.parentStore != nil {
		return d.parentStore.(*Directory).MetaData()
	}
	res = map[string]string{}
	dir, err := os.Open(d.Path)
	if err != nil {
		return
	}
	infos, err := dir.Readdir(0)
	dir.Close()
	if err != nil {
		return
	}
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		name := info.Name()
		if !(strings.HasPrefix(name, "._") && strings.HasSuffix(name, ".meta")) {
			continue
		}
		key, err := url.QueryUnescape(name)
		if err != nil {
			continue
		}
		key = strings.TrimSuffix(strings.TrimPrefix(key, "._"), ".meta")
		buf, err := ioutil.ReadFile(path.Join(d.Path, name))
		if err != nil {
			continue
		}
		val := strings.TrimSpace(string(buf))
		if val == "" {
			continue
		}
		res[key] = val
	}
	return res
}

func (d *Directory) SetMetaData(vals map[string]string) error {
	d.Lock()
	defer d.Unlock()
	if d.parentStore != nil {
		return d.parentStore.(*Directory).SetMetaData(vals)
	}
	written := map[string]struct{}{}
	for k, v := range vals {
		fileName := d.filename("._" + k + ".meta")
		if err := ioutil.WriteFile(fileName, []byte(v), 0644); err != nil {
			panic(err.Error())
		}
		written[path.Base(d.filename("._"+k+".meta"))] = struct{}{}
	}
	// Clean out the metadata values we no longer want
	dir, err := os.Open(d.Path)
	if err != nil {
		return err
	}
	infos, err := dir.Readdir(0)
	dir.Close()
	if err != nil {
		return err
	}
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		if _, ok := written[info.Name()]; ok {
			continue
		}
		os.Remove(path.Join(d.Path, info.Name()))
	}
	if n, ok := vals["Name"]; ok {
		d.name = n
	}
	return nil
}

func (f *Directory) Open(codec Codec) error {
	if f.Path == "" {
		return fmt.Errorf("Cannot store data at ''")
	}
	fullPath, err := filepath.Abs(filepath.Clean(f.Path))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return err
	}
	if codec == nil {
		codec = DefaultCodec
	}
	f.Codec = codec
	d, err := os.Open(fullPath)
	if err != nil {
		return err
	}
	infos, err := d.Readdir(0)
	d.Close()
	if err != nil {
		return err
	}
	f.opened = true
	for _, info := range infos {
		if info.IsDir() && info.Name() != "." && info.Name() != ".." {
			if _, err := f.MakeSub(info.Name()); err != nil {
				return err
			}
		}
	}
	md := f.MetaData()
	if n, ok := md["Name"]; ok {
		f.name = n
	}
	return nil
}

func (f *Directory) MakeSub(path string) (Store, error) {
	f.Lock()
	defer f.Unlock()
	f.panicIfClosed()
	if child, ok := f.subStores[path]; ok {
		return child, nil
	}
	child := &Directory{Path: filepath.Join(f.Path, path)}
	err := child.Open(f.Codec)
	if err != nil {
		return nil, err
	}
	addSub(f, child, path)
	return child, nil
}

func (f *Directory) Keys() ([]string, error) {
	f.panicIfClosed()
	d, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	infos, err := d.Readdir(0)
	d.Close()
	if err != nil {
		return nil, fmt.Errorf("dir keys: readdir error %#v", err)
	}
	res := []string{}
	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		name := info.Name()
		if !strings.HasSuffix(name, f.Ext()) {
			continue
		}
		n, err := url.QueryUnescape(strings.TrimSuffix(name, f.Ext()))
		if err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil
}

func (f *Directory) Load(key string, val interface{}) error {
	f.panicIfClosed()
	buf, err := ioutil.ReadFile(f.filename(key + f.Ext()))
	if err != nil {
		return err
	}
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

func (f *Directory) Save(key string, val interface{}) error {
	f.panicIfClosed()
	if f.ReadOnly() {
		return UnWritable(key)
	}
	buf, err := f.Encode(val)
	if err != nil {
		return err
	}
	return safeReplace(f.filename(key+f.Ext()), buf)
}

func (f *Directory) Remove(key string) error {
	f.panicIfClosed()
	if f.ReadOnly() {
		return UnWritable(key)
	}
	return os.Remove(f.filename(key + f.Ext()))
}
