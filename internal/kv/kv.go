package kv

import (
	"fmt"
	"path/filepath"

	"github.com/bull-cli/bull/internal/config"
	bolt "go.etcd.io/bbolt"
)

var defaultBucket = []byte("default")

func dbPath(name string) string {
	return filepath.Join(config.KVDir(), name+".db")
}

func openDB(name string) (*bolt.DB, error) {
	return bolt.Open(dbPath(name), 0600, nil)
}

func Put(dbName, bucket, key, value string) error {
	db, err := openDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	b := []byte(bucket)
	if bucket == "" {
		b = defaultBucket
	}
	return db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(b)
		if err != nil {
			return err
		}
		return bkt.Put([]byte(key), []byte(value))
	})
}

func Get(dbName, bucket, key string) (string, error) {
	db, err := openDB(dbName)
	if err != nil {
		return "", err
	}
	defer db.Close()

	b := []byte(bucket)
	if bucket == "" {
		b = defaultBucket
	}
	var val []byte
	err = db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(b)
		if bkt == nil {
			return fmt.Errorf("bucket %q not found", bucket)
		}
		v := bkt.Get([]byte(key))
		if v == nil {
			return fmt.Errorf("key %q not found", key)
		}
		val = make([]byte, len(v))
		copy(val, v)
		return nil
	})
	return string(val), err
}

func Del(dbName, bucket, key string) error {
	db, err := openDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	b := []byte(bucket)
	if bucket == "" {
		b = defaultBucket
	}
	return db.Update(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(b)
		if bkt == nil {
			return nil
		}
		return bkt.Delete([]byte(key))
	})
}

type KVPair struct {
	Key   string
	Value string
}

func List(dbName, bucket, prefix string) ([]KVPair, error) {
	db, err := openDB(dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	b := []byte(bucket)
	if bucket == "" {
		b = defaultBucket
	}
	var pairs []KVPair
	err = db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(b)
		if bkt == nil {
			return nil
		}
		c := bkt.Cursor()
		if prefix == "" {
			for k, v := c.First(); k != nil; k, v = c.Next() {
				pairs = append(pairs, KVPair{Key: string(k), Value: string(v)})
			}
		} else {
			pre := []byte(prefix)
			for k, v := c.Seek(pre); k != nil && len(k) >= len(pre) && string(k[:len(pre)]) == prefix; k, v = c.Next() {
				pairs = append(pairs, KVPair{Key: string(k), Value: string(v)})
			}
		}
		return nil
	})
	return pairs, err
}

func ListDBs() ([]string, error) {
	pattern := filepath.Join(config.KVDir(), "*.db")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, m := range matches {
		name := filepath.Base(m)
		names = append(names, name[:len(name)-3])
	}
	return names, nil
}
