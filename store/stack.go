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
	keys       map[string]int
}

func (s *StackedStore) Type() string {
	return "stacked"
}

func (s *StackedStore) Open(codec Codec) error {
	s.Codec = codec
	s.stores = []Store{}
	s.storeFlags = []layerFlags{}
	s.keys = map[string]int{}
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
	newLayerKeys []string
	err          error
	newSub       bool
	pushing      bool
	subTrackers  map[string]*pushTracker
}

func (pt *pushTracker) unlock() {
	for _, v := range pt.subTrackers {
		v.unlock()
	}
	pt.newLayer.RUnlock()
	if len(pt.stores) > 1 && pt.pushing {
		pt.newLayer.SetReadOnly()
	}
	pt.Unlock()
}

func (pt *pushTracker) push(kCBO, kCO bool) {
	pt.pushing = true
	for k, st := range pt.subTrackers {
		st.push(kCBO, kCO)
		if st.newSub {
			addSub(pt.StackedStore, st.StackedStore, k)
		}
	}
	newFlags := layerFlags{
		keysCannotBeOverridden: kCBO,
		keysCannotOverride:     kCO,
	}
	pt.storeFlags = append(pt.storeFlags, newFlags)
	pt.stores = append(pt.stores, pt.newLayer)
	for _, key := range pt.newLayerKeys {
		if _, ok := pt.keys[key]; !ok {
			pt.keys[key] = len(pt.stores) - 1
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
		subTrackers:  map[string]*pushTracker{},
	}
	res.newLayerKeys, res.err = layer.Keys()
	if res.err != nil {
		return
	}
	badKeys := []string{}
	for _, k := range res.newLayerKeys {
		i, ok := s.keys[k]
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
	if len(badKeys) != 0 {
		res.err = StackPushError(fmt.Sprintf("New layer violates key restrictions: %s", strings.Join(badKeys, "\n\t")))
		return
	}
	for k, v := range layer.Subs() {
		var subPT *pushTracker
		if subStore, ok := s.subStores[k]; !ok {
			newStore := &StackedStore{}
			newStore.Open(s.Codec)
			subPT = newStore.pushOK(v, kCBO, kCO)
			subPT.newSub = true
		} else {
			subPT = subStore.(*StackedStore).pushOK(v, kCBO, kCO)
		}
		res.subTrackers[k] = subPT
		if subPT.err != nil {
			res.err = subPT.err
			return
		}
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

func (s *StackedStore) MakeSub(st string) (Store, error) {
	s.Lock()
	defer s.Unlock()
	s.panicIfClosed()
	var mySub *StackedStore
	var err error
	if sub, ok := s.subStores[st]; ok {
		mySub = sub.(*StackedStore)
		mySub.Lock()
		defer mySub.Unlock()
	}
	sub := s.stores[0].GetSub(st)
	if sub != nil && mySub != nil {
		return mySub, nil
	}
	if sub == nil {
		sub, err = s.stores[0].MakeSub(st)
		if err != nil {
			return nil, err
		}
	}
	newSub := &StackedStore{}
	newSub.Open(s.Codec)
	err = newSub.Push(sub,
		s.storeFlags[0].keysCannotBeOverridden,
		s.storeFlags[0].keysCannotOverride)
	if err != nil {
		return nil, err
	}
	if mySub != nil {
		for i, sub := range mySub.stores {
			kCBO := mySub.storeFlags[i].keysCannotBeOverridden
			kCO := mySub.storeFlags[i].keysCannotOverride
			if err := newSub.Push(sub, kCBO, kCO); err != nil {
				return nil, err
			}
		}
		mySub.opened = false
	}
	addSub(s, newSub, st)
	return newSub, nil
}

func (s *StackedStore) Keys() ([]string, error) {
	s.RLock()
	defer s.RUnlock()
	vals := make([]string, 0, len(s.keys))
	for k := range s.keys {
		vals = append(vals, k)
	}
	return vals, nil
}

func (s *StackedStore) MetaFor(key string) map[string]string {
	s.RLock()
	defer s.RUnlock()
	idx, ok := s.keys[key]
	if !ok {
		return map[string]string{}
	}
	if ms, ok := s.stores[idx].(MetaSaver); ok {
		return ms.MetaData()
	}
	return map[string]string{}
}

func (s *StackedStore) Load(key string, val interface{}) error {
	s.RLock()
	defer s.RUnlock()
	idx, ok := s.keys[key]
	if !ok {
		return os.ErrNotExist
	}
	return s.stores[idx].Load(key, val)
}

type StackCannotOverride string

func (s StackCannotOverride) Error() string {
	return string(s)
}

type StackCannotBeOverridden string

func (s StackCannotBeOverridden) Error() string {
	return string(s)
}

func (s *StackedStore) Save(key string, val interface{}) error {
	s.RLock()
	defer s.RUnlock()
	idx, ok := s.keys[key]
	if ok && idx != 0 {
		// Key already exists.  Can it be overridden?
		if s.storeFlags[idx].keysCannotBeOverridden {
			return StackCannotBeOverridden(key)
		}
		if s.storeFlags[0].keysCannotOverride {
			return StackCannotOverride(key)
		}
	}
	err := s.stores[0].Save(key, val)
	if err == nil {
		s.keys[key] = 0
	}
	return err
}

func (s *StackedStore) Remove(key string) error {
	s.RLock()
	defer s.RUnlock()
	idx, ok := s.keys[key]
	if !ok {
		return os.ErrNotExist
	}
	if idx != 0 {
		return UnWritable(key)
	}
	err := s.stores[0].Remove(key)
	if err == nil {
		delete(s.keys, key)
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
