package kv

import (
	"fmt"
	"os"
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

func ListBuckets(dbName string) ([]string, error) {
	db, err := openDB(dbName)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var names []string
	err = db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			names = append(names, string(name))
			return nil
		})
	})
	return names, err
}

func Count(dbName, bucket string) (int, error) {
	db, err := openDB(dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	b := []byte(bucket)
	if bucket == "" {
		b = defaultBucket
	}
	var count int
	err = db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(b)
		if bkt == nil {
			return nil
		}
		count = bkt.Stats().KeyN
		return nil
	})
	return count, err
}

func ExportJSON(dbName, bucket string) ([]KVPair, error) {
	return List(dbName, bucket, "")
}

func ImportJSON(dbName, bucket string, pairs []KVPair) error {
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
		for _, p := range pairs {
			if err := bkt.Put([]byte(p.Key), []byte(p.Value)); err != nil {
				return err
			}
		}
		return nil
	})
}

func Exists(dbName, bucket, key string) (bool, error) {
	db, err := openDB(dbName)
	if err != nil {
		return false, err
	}
	defer db.Close()

	b := []byte(bucket)
	if bucket == "" {
		b = defaultBucket
	}
	var found bool
	err = db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(b)
		if bkt == nil {
			return nil
		}
		found = bkt.Get([]byte(key)) != nil
		return nil
	})
	return found, err
}

func Incr(dbName, bucket, key string, delta int64) (int64, error) {
	db, err := openDB(dbName)
	if err != nil {
		return 0, err
	}
	defer db.Close()

	b := []byte(bucket)
	if bucket == "" {
		b = defaultBucket
	}
	var newVal int64
	err = db.Update(func(tx *bolt.Tx) error {
		bkt, err := tx.CreateBucketIfNotExists(b)
		if err != nil {
			return err
		}
		v := bkt.Get([]byte(key))
		var cur int64
		if v != nil {
			fmt.Sscanf(string(v), "%d", &cur)
		}
		newVal = cur + delta
		return bkt.Put([]byte(key), []byte(fmt.Sprintf("%d", newVal)))
	})
	return newVal, err
}

func Scan(dbName, bucket, startKey, endKey string) ([]KVPair, error) {
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
		start := []byte(startKey)
		for k, v := c.Seek(start); k != nil; k, v = c.Next() {
			if len(endKey) > 0 && string(k) > endKey {
				break
			}
			pairs = append(pairs, KVPair{Key: string(k), Value: string(v)})
		}
		return nil
	})
	return pairs, err
}

func DropDB(dbName string) error {
	return os.Remove(dbPath(dbName))
}

func DropBucket(dbName, bucket string) error {
	db, err := openDB(dbName)
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket([]byte(bucket))
	})
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
