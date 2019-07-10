package store

import "os"

// MemoryStore provides an in-memory implementation of Store
// for testing purposes
type Memory struct {
	storeBase
	v    map[string]map[string][]byte
	meta map[string]string
}

func (m *Memory) Type() string {
	return "memory"
}

func (m *Memory) MetaData() map[string]string {
	m.RLock()
	defer m.RUnlock()
	res := map[string]string{}
	for k, v := range m.meta {
		res[k] = v
	}
	return res
}

func (m *Memory) SetMetaData(vals map[string]string) error {
	m.Lock()
	defer m.Unlock()
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
	m.v = map[string]map[string][]byte{}
	m.opened = true
	md := m.MetaData()
	if n, ok := md["Name"]; ok {
		m.name = n
	}
	return nil
}

func (m *Memory) Prefixes() ([]string, error) {
	m.RLock()
	defer m.RUnlock()
	m.panicIfClosed()
	res := []string{}
	for k := range m.v {
		res = append(res, k)
	}
	return res, nil
}

func (m *Memory) Keys(prefix string) ([]string, error) {
	m.RLock()
	m.panicIfClosed()
	defer m.RUnlock()
	res := []string{}
	vals, ok := m.v[prefix]
	if !ok {
		return res, nil
	}
	for k := range vals {
		res = append(res, k)
	}
	return res, nil
}

func (m *Memory) Exists(prefix, key string) bool {
	m.RLock()
	defer m.RUnlock()
	m.panicIfClosed()
	_, ok := m.v[prefix]
	if !ok {
		return ok
	}
	_, ok = m.v[prefix][key]
	return ok
}

func (m *Memory) Load(prefix, key string, val interface{}) error {
	m.RLock()
	defer m.RUnlock()
	m.panicIfClosed()
	if !m.Exists(prefix, key) {
		return os.ErrNotExist
	}
	if err := m.Decode(m.v[prefix][key], &val); err != nil {
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

func (m *Memory) Save(prefix, key string, val interface{}) error {
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
	if _, ok := m.v[prefix]; !ok {
		m.v[prefix] = map[string][]byte{}
	}
	m.v[prefix][key] = buf
	return nil
}

func (m *Memory) Remove(prefix, key string) error {
	m.Lock()
	defer m.Unlock()
	m.panicIfClosed()
	if _, ok := m.v[prefix]; !ok {
		return os.ErrNotExist
	}
	if _, ok := m.v[prefix][key]; !ok {
		return os.ErrNotExist
	}
	if m.readOnly {
		return UnWritable(key)
	}
	delete(m.v[prefix], key)
	return nil
}
