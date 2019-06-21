package store

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	consul "github.com/hashicorp/consul/api"
)

// Consul implements a Store that is backed by the Consul key/value store.
type Consul struct {
	storeBase
	source  string
	version string
	Client  *consul.Client

	BaseKey string
}

func (c *Consul) Open(codec Codec) error {
	if c.BaseKey == "" {
		return fmt.Errorf("Cannot store data at an empty location in the Consul KV store!")
	}
	if strings.HasPrefix(c.BaseKey, "/") {
		c.BaseKey = strings.TrimPrefix(c.BaseKey, "/")
	}
	if codec == nil {
		codec = DefaultCodec
	}
	c.Codec = codec
	if c.Client == nil {
		client, err := consul.NewClient(consul.DefaultConfig())
		if err != nil {
			return err
		}
		if info, err := client.Agent().Self(); err != nil {
			return err
		} else {
			c.source = fmt.Sprintf("consul: from %s", info["Config"]["NodeName"].(string))
		}
		c.Client = client
	}
	keys, qm, err := c.Client.KV().Keys(c.BaseKey, "", nil)
	if err != nil {
		return err
	}
	c.version = fmt.Sprintf("%d", qm.LastIndex)
	c.opened = true
	for i := range keys {
		if !strings.HasSuffix(keys[i], "/") {
			continue
		}
		subKey := strings.TrimSuffix(strings.TrimPrefix(keys[i], c.BaseKey+"/"), "/")
		if _, err := c.MakeSub(subKey); err != nil {
			return err
		}
	}
	c.closer = func() {
		c.Client = nil
	}
	md := c.MetaData()
	if n, ok := md["Name"]; ok {
		c.name = n
	}
	return nil
}

func (c *Consul) Type() string {
	return "consul"
}

func (b *Consul) MakeSub(prefix string) (Store, error) {
	b.Lock()
	defer b.Unlock()
	b.panicIfClosed()
	if res, ok := b.subStores[prefix]; ok {
		return res, nil
	}
	res := &Consul{Client: b.Client, BaseKey: filepath.Join(b.BaseKey, prefix)}
	err := res.Open(b.Codec)
	if err != nil {
		return nil, err
	}
	addSub(b, res, prefix)
	return res, nil
}

func (b *Consul) finalKey(k string) string {
	return path.Clean(path.Join(b.BaseKey, k))
}

func (b *Consul) Keys() ([]string, error) {
	b.panicIfClosed()
	keys, _, err := b.Client.KV().Keys(b.BaseKey, "", nil)
	if err != nil {
		return nil, err
	}
	res := []string{}
	for i := range keys {
		if strings.HasSuffix(keys[i], "/") {
			continue
		}
		res = append(res, strings.TrimPrefix(keys[i], b.BaseKey+"/"))
	}
	return res, nil
}

func (b *Consul) Load(key string, val interface{}) error {
	b.panicIfClosed()
	buf, _, err := b.Client.KV().Get(b.finalKey(key), nil)
	if buf == nil {
		return os.ErrNotExist
	}
	if err != nil {
		return err
	}
	if err := b.Decode(buf.Value, val); err != nil {
		return err
	}
	if ro, ok := val.(ReadOnlySetter); ok {
		ro.SetReadOnly(b.ReadOnly())
	}
	if bb, ok := val.(BundleSetter); ok {
		n := b.Name()
		if n != "" {
			bb.SetBundle(n)
		}
	}
	return nil
}

func (b *Consul) Save(key string, val interface{}) error {
	b.panicIfClosed()
	if b.ReadOnly() {
		return UnWritable(key)
	}
	buf, err := b.Encode(val)
	if err != nil {
		return err
	}
	kp := &consul.KVPair{Value: buf, Key: b.finalKey(key)}
	_, err = b.Client.KV().Put(kp, nil)
	return err
}

func (b *Consul) Remove(key string) error {
	b.panicIfClosed()
	if b.ReadOnly() {
		return UnWritable(key)
	}
	_, err := b.Client.KV().Delete(b.finalKey(key), nil)
	return err
}

func (b *Consul) MetaData() (res map[string]string) {
	if b.parentStore != nil {
		return b.parentStore.(*Consul).MetaData()
	}

	res = map[string]string{}
	b.Load("meta", &res)
	return res
}

func (b *Consul) SetMetaData(vals map[string]string) error {
	if b.parentStore != nil {
		return b.parentStore.(*Consul).SetMetaData(vals)
	}
	if n, ok := vals["Name"]; ok {
		b.name = n
	}
	return b.Save("meta", vals)
}
