package main

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

func initDB(file string) (*bolt.DB, error) {
	var db = new(bolt.DB)
	var err error
	for {
		db, err = bolt.Open(file, 0600, &bolt.Options{Timeout: time.Millisecond * 100})
		if err == bolt.ErrDatabaseOpen {
			continue
		}
		if err != nil {
			return nil, err
		}
		break
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("seen"))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("osmcache"))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("memo"))
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
