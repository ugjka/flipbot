package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"mvdan.cc/xurls/v2"
)

var (
	seenBucket     = []byte("seen")
	osmCacheBucket = []byte("osmcache")
	memoBucket     = []byte("memo")
	reminderBucket = []byte("reminder")
	logBucket      = []byte("log")
	indexBucket    = []byte("index")
	usersBucket    = []byte("users")
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
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(reminderBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(logBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(indexBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(usersBucket)
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

func setLogMSG(msg *Message) (err error) {
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(logBucket)
		id, _ := b.NextSequence()
		data, err := json.Marshal(msg)
		if err != nil {
			return err
		}
		err = b.Put(itob(id), data)
		if err != nil {
			return err
		}
		index := tx.Bucket(indexBucket)
		for _, v := range split(strings.ToLower(msg.Message)) {
			token, err := index.CreateBucketIfNotExists([]byte(v))
			if err != nil {
				return err
			}
			err = token.Put(itob(id), []byte(""))
			if err != nil {
				return err
			}
		}
		users := tx.Bucket(usersBucket)
		user, err := users.CreateBucketIfNotExists([]byte(strings.ToLower(msg.Nick)))
		if err != nil {
			return err
		}
		err = user.Put(itob(id), []byte(""))
		if err != nil {
			return err
		}
		return nil
	})
	return
}

func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btoi(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
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

func getReminder() (reminders ReminderItems, err error) {
	var tmp string
	err = db.View(func(tx *bolt.Tx) error {
		c := tx.Bucket(reminderBucket).Cursor()
		if k, v := c.First(); k != nil {
			tmp = string(k)
			rem, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return err
			}
			if rem.After(time.Now()) {
				return errNoReminder
			}
			return json.Unmarshal(v, &reminders)
		}
		return nil
	})
	if err != nil {
		return
	}
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(reminderBucket)
		return b.Delete([]byte(tmp))
	})
	return
}

func setReminder(target string, reminder ReminderItem) (err error) {
	reminders := make(ReminderItems, 0)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(reminderBucket)
		if v := b.Get([]byte(target)); v != nil {
			return json.Unmarshal(v, &reminders)
		}
		return nil
	})
	if err != nil {
		return
	}
	reminders = append(reminders, reminder)
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(reminderBucket)
		data, err := json.Marshal(reminders)
		if err != nil {
			return err
		}
		return b.Put([]byte(target), data)
	})
	return
}

var urls = xurls.Relaxed()
var word = regexp.MustCompile("(!?\\w+[\"']?\\w+)")

func split(in string) (out []string) {
	if urls.MatchString(in) {
		in = urls.ReplaceAllString(in, "")
	}
	in = strings.Replace(in, "’", "'", -1)
	return word.FindAllString(in, -1)
}

func search(in string, ignore string) (msgs []Message, err error) {
	var results [][]byte
	arr := split(in)
	store := make(map[uint64]int)
	err = db.View(func(tx *bolt.Tx) error {
		index := tx.Bucket(indexBucket)
		for _, v := range arr {
			if b := index.Bucket([]byte(v)); b != nil {
				b.ForEach(func(k, v []byte) error {
					store[btoi(k)]++
					return nil
				})
			} else {
				k, _ := index.Cursor().Seek([]byte(v))
				if k == nil {
					continue
				}
				if b := index.Bucket(k); b != nil {
					b.ForEach(func(k, v []byte) error {
						store[btoi(k)]++
						return nil
					})
				}
			}
		}
		depth := 0
		for _, v := range store {
			if v > depth {
				depth = v
			}
		}
		for k, v := range store {
			if v < depth {
				delete(store, k)
			}
		}
		for sk := range store {
			if b := index.Bucket([]byte(ignore)); b != nil {
				b.ForEach(func(k, v []byte) error {
					if sk == btoi(k) {
						delete(store, sk)
					}
					return nil
				})
			}
		}
		results = make([][]byte, len(store))
		i := 0
		for k := range store {
			results[i] = itob(k)
			i++
		}
		sort.Slice(results, func(i, j int) bool {
			res := bytes.Compare(results[i], results[j])
			if res == -1 {
				return true
			}
			return false
		})
		return nil
	})
	if len(results) == 0 {
		err = errNoResults
		return
	}
	err = db.View(func(tx *bolt.Tx) error {
		msgs = make([]Message, 0, 5)
		for i, j := len(results)-1, 0; i >= 0 && j < 5; i, j = i-1, j+1 {
			msg := Message{}
			err := json.Unmarshal(tx.Bucket(logBucket).Get(results[i]), &msg)
			if err != nil {
				return err
			}
			msgs = append([]Message{msg}, msgs...)
		}
		return nil
	})
	return
}

func userTail(nick string, ignore string, amount int) (msgs []Message, err error) {
	msgs = make([]Message, 0, amount)
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(usersBucket).Bucket([]byte(nick))
		if b == nil {
			return errNotSeen

		}
		c := b.Cursor()
		i := 0
		log := tx.Bucket(logBucket)
		for k, _ := c.Last(); k != nil && i < amount; k, _ = c.Prev() {
			if v := log.Get(k); v != nil {
				msg := Message{}
				err := json.Unmarshal(v, &msg)
				if strings.Contains(msg.Message, ignore) {
					continue
				}
				if err != nil {
					return err
				}
				msgs = append([]Message{msg}, msgs...)
				i++
			}
		}
		return nil
	})
	if len(msgs) == 0 {
		err = errNoResults
	}
	return
}
