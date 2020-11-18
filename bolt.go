package main

import (
	"bytes"
	"fmt"

	"github.com/boltdb/bolt"
)

var (
	osmCacheBucket = []byte("osmcache")
)

func initDB(file string) (*bolt.DB, error) {
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(osmCacheBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func setOSMCache(url string, data []byte) error {
	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(osmCacheBucket)
		return b.Put([]byte(url), data)
	})
}

func getOSMCache(url string) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(osmCacheBucket)
		data := b.Get([]byte(url))
		if data == nil {
			return errNotInCache
		}
		_, err := buf.Write(data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
