package store

import "os"

// MemoryStore provides an in-memory implementation of Store
// for testing purposes
type Memory struct {
	storeBase
	v    map[string][]byte
	meta map[string]string
}

func (m *Memory) Type() string {
	return "memory"
}

func (m *Memory) MetaData() map[string]string {
	m.RLock()
	defer m.RUnlock()
	if m.parentStore != nil {
		return m.parentStore.(*Memory).MetaData()
	}
	res := map[string]string{}
	for k, v := range m.meta {
		res[k] = v
	}
	return res
}

func (m *Memory) SetMetaData(vals map[string]string) error {
	m.Lock()
	defer m.Unlock()
	if m.parentStore != nil {
		return m.parentStore.(*Memory).SetMetaData(vals)
	}
	m.meta = map[string]string{}
	for k, v := range vals {
		m.meta[k] = v
	}
	if n, ok := vals["Name"]; ok {
		m.name = n
	}
	return nil
}

func (m *Memory) Open(codec Codec) error {
	if codec == nil {
		codec = DefaultCodec
	}
	m.Codec = codec
	m.closer = func() {
		m.v = nil
	}
	m.v = map[string][]byte{}
	m.opened = true
	md := m.MetaData()
	if n, ok := md["Name"]; ok {
		m.name = n
	}
	return nil
}

func (m *Memory) MakeSub(loc string) (Store, error) {
	m.Lock()
	defer m.Unlock()
	m.panicIfClosed()
	if res, ok := m.subStores[loc]; ok {
		return res, nil
	}
	res := &Memory{}
	res.Open(m.Codec)
	addSub(m, res, loc)
	return res, nil
}

func (m *Memory) Keys() ([]string, error) {
	m.RLock()
	m.panicIfClosed()
	res := make([]string, 0, len(m.v))
	for k := range m.v {
		res = append(res, k)
	}
	m.RUnlock()
	return res, nil
}

func (m *Memory) Load(key string, val interface{}) error {
	m.RLock()
	m.panicIfClosed()
	v, ok := m.v[key]
	m.RUnlock()
	if !ok {
		return os.ErrNotExist
	}
	if err := m.Decode(v, val); err != nil {
		return err
	}
	if ro, ok := val.(ReadOnlySetter); ok {
		ro.SetReadOnly(m.ReadOnly())
	}
	if bb, ok := val.(BundleSetter); ok {
		n := m.Name()
		if n != "" {
			bb.SetBundle(n)
		}
	}
	return nil
}

func (m *Memory) Save(key string, val interface{}) error {
	m.Lock()
	defer m.Unlock()
	m.panicIfClosed()
	if m.readOnly {
		return UnWritable(key)
	}
	buf, err := m.Encode(val)
	if err != nil {
		return err
	}
	m.v[key] = buf
	return nil
}

func (m *Memory) Remove(key string) error {
	m.Lock()
	defer m.Unlock()
	m.panicIfClosed()
	_, ok := m.v[key]
	if ok {
		if m.readOnly {
			return UnWritable(key)
		}
		delete(m.v, key)
		return nil
	}
	return os.ErrNotExist
}
