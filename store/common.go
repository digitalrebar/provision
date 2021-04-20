package store

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"path"
	"sync"
)

func safeReplace(name string, contents []byte) error {
	if err := os.MkdirAll(path.Dir(name), 0700); err != nil {
		return err
	}
	tmpName := path.Join(path.Dir(name), ".new."+path.Base(name))
	f, err := os.OpenFile(tmpName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0640)
	if err != nil {
		return err
	}
	func() {
		defer f.Close()
		if _, err = f.Write(contents); err == nil {
			err = f.Sync()
		}
	}()
	if err != nil {
		return err
	}
	return os.Rename(tmpName, name)
}

// Open a store via URI style locator. Locators have the following formats:
//
// storeType:path?codec=codecType&ro=false&option=foo for stores
// that always refer to something local, and
//
// storeType://host:port/path?codec=codecType&ro=false&option=foo for stores
// that need to talk over the network.
//
// All store types take codec and ro as optional parameters
//
// The following storeTypes are known:
//   * file, in which path refers to a single local file.
//   * directory, in which path refers to a top-level directory
//   * bolt, in which path refers to the directory where the Bolt database
//     is located.  bolt also takes an optional bucket parameter to specify the
//     top-level bucket data is stored in.
//   * memory, in which path does not mean anything.
//
func Open(locator string) (Store, error) {
	uri, err := url.Parse(locator)
	if err != nil {
		return nil, err
	}
	params := uri.Query()
	codec := DefaultCodec
	readOnly := false
	codecParam := params.Get("codec")
	switch codecParam {
	case "yaml":
		codec = YamlCodec
	case "json":
		codec = JsonCodec
	case "", "default":
		codec = DefaultCodec
	default:
		return nil, fmt.Errorf("Unknown codec %s", codecParam)
	}
	roParam := params.Get("ro")
	switch roParam {
	case "true", "yes", "1":
		readOnly = true
	case "false", "no", "0", "":
		readOnly = false
	default:
		return nil, fmt.Errorf("Unknown ro value %s. Try true or false", roParam)
	}
	var res Store
	urlPath := uri.Opaque
	if urlPath == "" {
		urlPath = uri.Path
	}
	switch uri.Scheme {
	case "stack":
		res = &StackedStore{}
	case "file":
		res = &File{Path: urlPath}
	case "directory":
		res = &Directory{Path: urlPath}
	case "memory":
		res = &Memory{}
	}
	if res == nil {
		return nil, fmt.Errorf("Unknown schema type: %s", uri.Scheme)
	}
	if err := res.Open(codec); err != nil {
		return nil, err
	}
	if readOnly {
		res.SetReadOnly()
	}
	return res, nil
}

// Store provides an interface for some very basic key/value
// storage needs.  Each Store (including ones created with MakeSub()
// should operate as seperate, flat key/value stores.
type Store interface {
	sync.Locker
	RLock()
	RUnlock()
	// Open opens the store for use.
	Open(Codec) error
	// GetCodec returns the codec that the open store uses for marshalling and unmarshalling data
	GetCodec() Codec
	// Keys returns the list of keys that this store has in no
	// particular order.
	Keys(string) ([]string, error)
	// Subs returns a map all of the substores for this store.
	Prefixes() ([]string, error)
	// Test to see if a given entry exists
	Exists(string, string) bool
	// Load the data for a particular key
	Load(string, string, interface{}) error
	// Save data for a key
	Save(string, string, interface{}) error
	// Remove a key/value pair.
	Remove(string, string) error
	// ReadOnly returns whether a store is set to be read-only.
	ReadOnly() bool
	// SetReadOnly sets the store into read-only mode.  This is a
	// one-way operation -- once a store is set to read-only, it
	// cannot be changed back to read-write while the store is open.
	SetReadOnly() bool
	// Close closes the store.  Attempting to perfrom operations on
	// a closed store will panic.
	Close()
	// Closed returns whether or not a store is Closed
	Closed() bool
	// Type is the type of Store this is.
	Type() string
	// Name is the name of Store this is.
	Name() string
}

// MetaSaver is a Store that is capable of of recording
// metadata about itself.
type MetaSaver interface {
	Store
	MetaData() map[string]string
	SetMetaData(map[string]string) error
}

// Copy copies all of the contents from src to dest, including substores and
// metadata.  If dst starts out empty, then dst will wind up being a clone of src.
func Copy(dst, src Store) error {
	src.RLock()
	defer src.RUnlock()
	dmeta, dok := dst.(MetaSaver)
	smeta, sok := src.(MetaSaver)
	if dok && sok {
		if err := dmeta.SetMetaData(smeta.MetaData()); err != nil {
			return err
		}
	}
	prefixes, err := src.Prefixes()
	if err != nil {
		return err
	}
	for _, prefix := range prefixes {
		keys, err := src.Keys(prefix)
		if err != nil {
			return err
		}
		for _, key := range keys {
			var val interface{}
			if err := src.Load(prefix, key, &val); err != nil {
				return err
			}
			if err := dst.Save(prefix, key, val); err != nil {
				return err
			}
		}
	}
	return nil
}

type forceCloser interface {
	forceClose()
}

// NotFound is the "key not found" error type.
type NotFound string

func (n NotFound) Error() string {
	return fmt.Sprintf("key %s: not found", string(n))
}

type UnWritable string

func (u UnWritable) Error() string {
	return fmt.Sprintf("readonly: %s", string(u))
}

type storeBase struct {
	sync.RWMutex
	Codec
	readOnly    bool
	opened      bool
	subStores   map[string]Store
	parentStore Store
	closer      func()
	name        string
}

func (s *storeBase) Name() string {
	if s.parentStore != nil {
		return s.parentStore.Name()
	}
	return s.name
}

func (s *storeBase) Close() {
	s.Lock()
	defer s.Unlock()
	if !s.opened {
		return
	}
	if s.closer != nil {
		s.closer()
	}
	s.opened = false
}

func (s *storeBase) GetCodec() Codec {
	return s.Codec
}

func (s *storeBase) panicIfClosed() {
	if !s.opened {
		panic("Operation on closed store")
	}
}

func (s *storeBase) ReadOnly() bool {
	s.RLock()
	defer s.RUnlock()
	s.panicIfClosed()
	return s.readOnly
}

func (s *storeBase) SetReadOnly() bool {
	s.Lock()
	defer s.Unlock()
	s.panicIfClosed()
	if s.readOnly {
		return false
	}
	s.readOnly = true
	for _, sub := range s.subStores {
		sub.SetReadOnly()
	}
	return true
}

func (s *storeBase) Closed() bool {
	return !s.opened
}

func remarshal(src, dest interface{}) error {
	buf, err := json.Marshal(src)
	if err == nil {
		err = json.Unmarshal(buf, dest)
	}
	return err
}
