package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/boltdb/bolt"
)

var (
	seenBucket     = []byte("seen")
	osmCacheBucket = []byte("osmcache")
	memoBucket     = []byte("memo")
)

func initDB(file string) (*bolt.DB, error) {
	db, err := bolt.Open(file, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(seenBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
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
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(memoBucket)
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

func getSeen(nick string) (seen Seen, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(seenBucket)
		res := b.Get([]byte(nick))
		if res == nil {
			return errNotSeen
		}
		return json.Unmarshal(res, &seen)
	})
	return
}

func getSeenPrefix(prefix string) (nick string, seen Seen, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(seenBucket).Cursor()
		if k, _ := c.Seek([]byte(prefix)); k == nil || !bytes.HasPrefix(k, []byte(prefix)) {
			return errNotSeen
		}
		var tmpSeen = Seen{}
		for k, v := c.Seek([]byte(prefix)); k != nil && bytes.HasPrefix(k, []byte(prefix)); k, v = c.Next() {
			err := json.Unmarshal(v, &tmpSeen)
			if err != nil {
				return err
			}
			if tmpSeen.Seen.After(seen.Seen) {
				seen = tmpSeen
				nick = string(k)
			}
		}
		return nil
	})
	return
}

func setSeen(nick string, seen *Seen) error {
	err := db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(seenBucket)
		data, err := json.Marshal(seen)
		if err != nil {
			return err
		}
		return b.Put([]byte(nick), data)
	})
	return err
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

func setMemo(nick string, memo memoItem) (err error) {
	items := make(memos, 0)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(memoBucket)
		if v := b.Get([]byte(nick)); v != nil {
			return json.Unmarshal(v, &items)
		}
		return nil
	})
	if err != nil {
		return
	}
	items = append(items, memo)
	return db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(memoBucket)
		data, err := json.Marshal(items)
		if err != nil {
			return err
		}
		return b.Put([]byte(nick), data)
	})
}

func getMemo(nick string) (items memos, err error) {
	items = make(memos, 0)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(memoBucket)
		if v := b.Get([]byte(nick)); v != nil {
			return json.Unmarshal(v, &items)
		}
		return errNoMemo
	})
	if err != nil {
		return
	}
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(memoBucket)
		return b.Delete([]byte(nick))
	})
	return
}
