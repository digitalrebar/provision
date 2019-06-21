package store

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"sync"
)

func safeReplace(name string, contents []byte) error {
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
	path := uri.Opaque
	if path == "" {
		path = uri.Path
	}
	switch uri.Scheme {
	case "stack":
		res = &StackedStore{}
	case "file":
		res = &File{Path: path}
	case "directory":
		res = &Directory{Path: path}
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
	// GetSub fetches an already-existing substore.  nil means there is no such substore.
	GetSub(string) Store
	// MakeSub returns a Store that is subordinate to this one.
	// What exactly that means depends on the simplestore in question,
	// but it should wind up sharing the same backing store (directory,
	// database, etcd cluster, whatever)
	MakeSub(string) (Store, error)
	// Parent fetches the parent of this store, if any.
	Parent() Store
	// Keys returns the list of keys that this store has in no
	// particular order.
	Keys() ([]string, error)
	// Subs returns a map all of the substores for this store.
	Subs() map[string]Store
	// Load the data for a particular key
	Load(string, interface{}) error
	// Save data for a key
	Save(string, interface{}) error
	// Remove a key/value pair.
	Remove(string) error
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
	keys, err := src.Keys()
	if err != nil {
		return err
	}
	for _, key := range keys {
		var val interface{}
		if err := src.Load(key, &val); err != nil {
			return err
		}
		if err := dst.Save(key, val); err != nil {
			return err
		}
	}
	for k, sub := range src.Subs() {
		subDst, err := dst.MakeSub(k)
		if err != nil {
			return err
		}
		if err := Copy(subDst, sub); err != nil {
			return err
		}
	}
	return nil
}

type parentSetter interface {
	setParent(Store)
}

type childSetter interface {
	addChild(string, Store)
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

func (s *storeBase) forceClose() {
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

func (s *storeBase) Close() {
	s.Lock()
	if s.parentStore == nil {
		s.Unlock()
		s.forceClose()
		for _, sub := range s.subStores {
			sub.(forceCloser).forceClose()
		}
		return
	}
	parent := s.parentStore
	s.Unlock()
	parent.Close()
	return
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

func (s *storeBase) GetSub(name string) Store {
	s.RLock()
	defer s.RUnlock()
	s.panicIfClosed()
	if s.subStores == nil {
		return nil
	}
	return s.subStores[name]
}

func (s *storeBase) Subs() map[string]Store {
	s.RLock()
	defer s.RUnlock()
	s.panicIfClosed()
	res := map[string]Store{}
	for k, v := range s.subStores {
		res[k] = v
	}
	return res
}

func (s *storeBase) Parent() Store {
	s.RLock()
	defer s.RUnlock()
	s.panicIfClosed()
	return s.parentStore.(Store)
}

func (s *storeBase) Closed() bool {
	return !s.opened
}

func (s *storeBase) setParent(p Store) {
	s.parentStore = p
}

func (s *storeBase) addChild(name string, c Store) {
	if s.subStores == nil {
		s.subStores = map[string]Store{}
	}
	s.subStores[name] = c
}

func addSub(parent, child Store, name string) {
	parent.(childSetter).addChild(name, child)
	child.(parentSetter).setParent(parent)
}
