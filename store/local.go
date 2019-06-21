package store

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
)

type Bolt struct {
	storeBase
	Path   string
	db     *bolt.DB
	Bucket []byte
}

func (b *Bolt) Type() string {
	return "bolt"
}

func (b *Bolt) MetaData() map[string]string {
	res := map[string]string{}
	b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("$metadata"))
		if bucket == nil {
			return nil
		}
		bucket.ForEach(func(k, v []byte) error {
			if v == nil {
				return nil
			}
			res[string(k)] = string(v)
			return nil
		})
		return nil
	})
	return res
}

func (b *Bolt) SetMetaData(vals map[string]string) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket([]byte("$metadata"))
		bucket, err := tx.CreateBucket([]byte("$metadata"))
		if err != nil {
			return err
		}
		for k, v := range vals {
			if err := bucket.Put([]byte(k), []byte(v)); err != nil {
				return err
			}
		}
		if n, ok := vals["Name"]; ok {
			b.name = n
		}
		return nil
	})
}

func (b *Bolt) getBucket(tx *bolt.Tx) (res *bolt.Bucket) {
	for _, part := range bytes.Split(b.Bucket, []byte("/")) {
		if res == nil {
			res = tx.Bucket(part)
		} else {
			res = res.Bucket(part)
		}
		if res == nil {
			panic(fmt.Sprintf("Bucket %s does not exist", string(b.Bucket)))
		}
	}
	return
}

func (b *Bolt) createBucket() error {
	return b.db.Update(func(tx *bolt.Tx) error {
		var bukkit *bolt.Bucket
		var err error
		for _, part := range bytes.Split(b.Bucket, []byte("/")) {
			if bukkit == nil {
				bukkit, err = tx.CreateBucketIfNotExists(part)
			} else {
				bukkit, err = bukkit.CreateBucketIfNotExists(part)
			}
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (b *Bolt) MakeSub(loc string) (Store, error) {
	b.Lock()
	defer b.Unlock()
	b.panicIfClosed()
	if res, ok := b.subStores[loc]; ok {
		return res, nil
	}
	res := &Bolt{db: b.db, Bucket: bytes.Join([][]byte{b.Bucket, []byte(loc)}, []byte("/"))}
	if err := res.Open(b.Codec); err != nil {
		return nil, err
	}
	res.closer = func() {
		res.db = nil
	}
	addSub(b, res, loc)
	return res, nil
}

func (l *Bolt) loadSubs() error {
	if err := l.createBucket(); err != nil {
		return err
	}
	subs := [][]byte{}
	err := l.db.View(func(tx *bolt.Tx) error {
		bucket := l.getBucket(tx)
		bucket.ForEach(func(k, v []byte) error {
			if v == nil {
				subs = append(subs, k)
			}
			return nil
		})
		return nil
	})
	if err != nil {
		return err
	}

	for _, sub := range subs {
		if _, err := l.MakeSub(string(sub)); err != nil {
			return err
		}
	}
	return nil
}

func (b *Bolt) Open(codec Codec) error {
	if b.Bucket == nil {
		b.Bucket = []byte(`Default`)
	}
	if codec == nil {
		codec = DefaultCodec
	}
	b.Codec = codec
	if b.db == nil {
		if b.Path == "" {
			return fmt.Errorf("Cannot store data in ''")
		}
		finalLoc := filepath.Clean(b.Path)
		if err := os.MkdirAll(finalLoc, 0755); err != nil {
			return err
		}
		db, err := bolt.Open(filepath.Join(finalLoc, "bolt.db"), 0600, nil)
		if err != nil {
			return err
		}
		b.db = db
	}
	b.opened = true
	if err := b.loadSubs(); err != nil {
		return err
	}

	b.closer = func() {
		b.db.Close()
		b.db = nil
	}
	md := b.MetaData()
	if n, ok := md["Name"]; ok {
		b.name = n
	}
	return nil
}

func (b *Bolt) Keys() ([]string, error) {
	b.panicIfClosed()
	res := []string{}
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := b.getBucket(tx)
		bucket.ForEach(func(k, v []byte) error {
			if v != nil {
				res = append(res, string(k))
			}
			return nil
		})
		return nil
	})
	return res, err
}

func (b *Bolt) Load(key string, val interface{}) error {
	b.panicIfClosed()
	var res []byte
	err := b.db.View(func(tx *bolt.Tx) error {
		bucket := b.getBucket(tx)
		res = bucket.Get([]byte(key))
		if res == nil {
			return os.ErrNotExist
		}
		return nil
	})
	if err != nil {
		return err
	}
	if err := b.Decode(res, val); err != nil {
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

func (b *Bolt) Save(key string, val interface{}) error {
	b.panicIfClosed()
	if b.ReadOnly() {
		return UnWritable(key)
	}
	buf, err := b.Encode(val)
	if err != nil {
		return err
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := b.getBucket(tx)
		return bucket.Put([]byte(key), buf)
	})
}

func (b *Bolt) Remove(key string) error {
	b.panicIfClosed()
	if b.ReadOnly() {
		return UnWritable(key)
	}
	return b.db.Update(func(tx *bolt.Tx) error {
		bucket := b.getBucket(tx)
		if res := bucket.Get([]byte(key)); res == nil {
			return os.ErrNotExist
		}
		return bucket.Delete([]byte(key))
	})
}
