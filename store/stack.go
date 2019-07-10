package store

import (
	"fmt"
	"os"
	"strings"
)

type layerFlags struct {
	// Test to see if layers from n-1..0 have the same key.  If they do,
	// then the keys will override this one, violating the stack sanity
	// checking rules.
	keysCannotBeOverridden bool
	// Test to see if layers from n+1..len(stack) have this key.  If
	// they do, then this key will override that key, violating the
	// stack sanity checking rules.
	keysCannotOverride bool
}

// StackedStore is a store that represents the combination of several
// stores stacked together.  The first store in the stack is the only
// one that is writable, and the rest are set as read-only.
// StackedStores are initally created empty.
type StackedStore struct {
	storeBase
	stores     []Store
	storeFlags []layerFlags
	keys       map[string]map[string]int
}

func (s *StackedStore) Type() string {
	return "stacked"
}

func (s *StackedStore) Open(codec Codec) error {
	s.Codec = codec
	s.stores = []Store{}
	s.storeFlags = []layerFlags{}
	s.keys = map[string]map[string]int{}
	s.opened = true
	s.closer = func() {
		for _, item := range s.stores {
			item.Close()
		}
	}
	return nil
}

type pushTracker struct {
	*StackedStore
	newLayer     Store
	newLayerKeys map[string][]string
	err          error
	pushing      bool
}

func (pt *pushTracker) unlock() {
	pt.newLayer.RUnlock()
	if len(pt.stores) > 1 && pt.pushing {
		pt.newLayer.SetReadOnly()
	}
	pt.Unlock()
}

func (pt *pushTracker) push(kCBO, kCO bool) {
	pt.pushing = true
	newFlags := layerFlags{
		keysCannotBeOverridden: kCBO,
		keysCannotOverride:     kCO,
	}
	pt.storeFlags = append(pt.storeFlags, newFlags)
	pt.stores = append(pt.stores, pt.newLayer)
	for prefix, keys := range pt.newLayerKeys {
		if _, ok := pt.keys[prefix]; !ok {
			pt.keys[prefix] = map[string]int{}
		}
		for _, key := range keys {
			if _, ok := pt.keys[prefix][key]; !ok {
				pt.keys[prefix][key] = len(pt.stores) - 1
			}
		}
	}
}

type StackPushError string

func (s StackPushError) Error() string {
	return string(s)
}

func (s *StackedStore) pushOK(layer Store, kCBO, kCO bool) (res *pushTracker) {
	s.Lock()
	layer.RLock()
	s.panicIfClosed()
	if layer.Closed() {
		panic("Cannot push a closed store")
	}
	res = &pushTracker{
		StackedStore: s,
		newLayer:     layer,
		newLayerKeys: map[string][]string{},
	}
	prefixes, err := layer.Prefixes()
	if err != nil {
		res.err = err
		return
	}
	badKeys := []string{}
	for _, prefix := range prefixes {
		res.newLayerKeys[prefix], res.err = layer.Keys(prefix)
		if res.err != nil {
			return
		}
		for _, k := range res.newLayerKeys[prefix] {
			i, ok := s.keys[prefix][k]
			if !ok {
				// New key.  Cannot be overridden, and nothing else would override it that should not.
				continue
			}
			if kCBO {
				badKeys = append(badKeys,
					fmt.Sprintf("keysCannotBeOverridden: %s is already in layer %d", k, i))
			}
			if s.storeFlags[i].keysCannotOverride {
				badKeys = append(badKeys,
					fmt.Sprintf("keysCannotOverride: %s would be overridden by layer %d", k, i))
			}
		}
	}
	if len(badKeys) != 0 {
		res.err = StackPushError(fmt.Sprintf("New layer violates key restrictions: %s", strings.Join(badKeys, "\n\t")))
	}
	return
}

// Push adds a Store to the stack of stores in this stack.  Any Store
// but the inital one will be marked as read-only.  Either the Push
// call succeeds, or nothing about any of the Stores that are part of
// the Push operation change and the error contains details about what
// went wrong.
func (s *StackedStore) Push(layer Store, keysCannotBeOverridden, keysCannotOverride bool) error {
	tracker := s.pushOK(layer, keysCannotBeOverridden, keysCannotOverride)
	defer tracker.unlock()
	if tracker.err != nil {
		return tracker.err
	}
	tracker.push(keysCannotBeOverridden, keysCannotOverride)
	return nil
}

func (s *StackedStore) Layers() []Store {
	s.Lock()
	defer s.Unlock()
	res := make([]Store, len(s.stores))
	copy(res, s.stores)
	return res
}

func (s *StackedStore) Prefixes() ([]string, error) {
	s.RLock()
	defer s.RUnlock()
	vals := make([]string, 0, len(s.keys))
	for k := range s.keys {
		vals = append(vals, k)
	}
	return vals, nil
}

func (s *StackedStore) Keys(prefix string) ([]string, error) {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.keys[prefix]; !ok {
		return []string{}, nil
	}
	vals := make([]string, 0, len(s.keys[prefix]))
	for k := range s.keys[prefix] {
		vals = append(vals, k)
	}
	return vals, nil
}

func (s *StackedStore) MetaFor(prefix, key string) map[string]string {
	s.RLock()
	defer s.RUnlock()
	idx, ok := s.keys[prefix][key]
	if !ok {
		return map[string]string{}
	}
	if ms, ok := s.stores[idx].(MetaSaver); ok {
		return ms.MetaData()
	}
	return map[string]string{}
}

func (s *StackedStore) ItemReadOnly(prefix, key string) (bool, bool) {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.keys[prefix]; !ok {
		s.keys[prefix] = map[string]int{}
	}
	i, ok := s.keys[prefix][key]
	return i != 0, ok
}

func (s *StackedStore) Exists(prefix, key string) bool {
	_, res := s.ItemReadOnly(prefix, key)
	return res
}

func (s *StackedStore) Load(prefix, key string, val interface{}) error {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.keys[prefix]; !ok {
		s.keys[prefix] = map[string]int{}
	}
	idx, ok := s.keys[prefix][key]
	if !ok {
		return os.ErrNotExist
	}
	return s.stores[idx].Load(prefix, key, val)
}

type StackCannotOverride string

func (s StackCannotOverride) Error() string {
	return string(s)
}

type StackCannotBeOverridden string

func (s StackCannotBeOverridden) Error() string {
	return string(s)
}

func (s *StackedStore) Save(prefix, key string, val interface{}) error {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.keys[prefix]; !ok {
		s.keys[prefix] = map[string]int{}
	}
	idx, ok := s.keys[prefix][key]
	if ok && idx != 0 {
		// Key already exists.  Can it be overridden?
		if s.storeFlags[idx].keysCannotBeOverridden {
			return StackCannotBeOverridden(key)
		}
		if s.storeFlags[0].keysCannotOverride {
			return StackCannotOverride(key)
		}
	}
	err := s.stores[0].Save(prefix, key, val)
	if err == nil {
		s.keys[prefix][key] = 0
	}
	return err
}

func (s *StackedStore) Remove(prefix, key string) error {
	s.RLock()
	defer s.RUnlock()
	if _, ok := s.keys[prefix]; !ok {
		return os.ErrNotExist
	}
	idx, ok := s.keys[prefix][key]
	if !ok {
		return os.ErrNotExist
	}
	if idx != 0 {
		return UnWritable(key)
	}
	err := s.stores[0].Remove(prefix, key)
	if err == nil {
		delete(s.keys[prefix], key)
	}
	return err
}

func (s *StackedStore) ReadOnly() bool {
	s.RLock()
	defer s.RUnlock()
	return s.stores[0].ReadOnly()
}

func (s *StackedStore) SetReadOnly() bool {
	s.RLock()
	defer s.RUnlock()
	return s.stores[0].SetReadOnly()
}
