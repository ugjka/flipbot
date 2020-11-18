package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/dustin/go-humanize"

	"github.com/boltdb/bolt"
	"github.com/hako/durafmt"
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
	rankBucket     = []byte("ranks")
	hostmaskBucket = []byte("hostmask")
	quietBucket    = []byte("quiet")
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
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(rankBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return err
	})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(hostmaskBucket)
		if err != nil {
			return fmt.Errorf("create bucket: %v", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(quietBucket)
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

type dead struct {
	times []struct {
		duration time.Duration
		posts    int
	}
	last struct {
		seen time.Time
		nick string
	}
}

func (d dead) String() string {
	str := ""
	str += fmt.Sprintf("Last activity by %s %s. ", d.last.nick[:len(d.last.nick)-1]+"*", humanize.Time(d.last.seen))
	if len(d.times) != 0 {
		str += "Total: "
	}
	for _, v := range d.times {
		if v.posts == 1 {
			str += fmt.Sprintf("[%d post in last %s] ", v.posts, roundDuration(durafmt.Parse(v.duration).String()))
		} else {
			str += fmt.Sprintf("[%d posts in last %s] ", v.posts, roundDuration(durafmt.Parse(v.duration).String()))
		}

	}
	return str
}

func getDead(dur ...time.Duration) (d dead, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(logBucket)
		c := b.Cursor()
		counter := 0
		now := time.Now().UTC()
		msg := Message{}
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			err := json.Unmarshal(v, &msg)
			if err != nil {
				return err
			}
			if strings.HasPrefix(msg.Message, "!") || strings.HasPrefix(msg.Message, "$") {
				continue
			}
			if d.last.nick == "" {
				d.last.nick = msg.Nick
				d.last.seen = msg.Time
			}
			if len(dur) == 0 {
				break
			}
			if now.Add(-dur[0]).After(msg.Time) {
				d.times = append(d.times, struct {
					duration time.Duration
					posts    int
				}{
					duration: dur[0],
					posts:    counter,
				})
				dur = dur[1:]
				continue
			}
			counter++
		}
		return nil
	})
	return d, err
}

type recent []struct {
	nick string
	time time.Time
}

func (r recent) String() string {
	str := "Recent activity: "
	for _, v := range r {
		str += fmt.Sprintf("[%s: %s] ", v.nick[:len(v.nick)-1]+"*", humanize.Time(v.time))
	}
	return str
}

func getRecent(items int) (r recent, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(logBucket)
		c := b.Cursor()
		msg := Message{}
		memory := make(map[string]struct{})
		nick := ""
		for k, v := c.Last(); k != nil; k, v = c.Prev() {
			if items == 0 {
				break
			}
			err := json.Unmarshal(v, &msg)
			if err != nil {
				return err
			}
			nick = strings.ToLower(msg.Nick)
			if strings.HasPrefix(msg.Message, "!") || strings.HasPrefix(msg.Message, "$") {
				continue
			}
			if _, ok := memory[nick]; ok {
				continue
			}
			memory[nick] = struct{}{}
			r = append(r, struct {
				nick string
				time time.Time
			}{
				nick: nick,
				time: msg.Time,
			})
			items--
		}
		return nil
	})
	return r, err
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

func setvote(value int64, word string) (votes float64, err error) {
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(rankBucket)
		wb, err := b.CreateBucketIfNotExists([]byte(word))
		if err != nil {
			return err
		}
		buf := make([]byte, binary.MaxVarintLen16)
		n := binary.PutVarint(buf, value)
		now, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		err = wb.Put([]byte(now.Format(time.RFC3339)), buf[:n])
		if err != nil {
			return err
		}
		c := wb.Cursor()
		week := (time.Hour * 24 * 7).Seconds()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			voteDate, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return err
			}
			if voteDate.Before(now.Add(time.Hour * -24 * 7)) {
				err = c.Delete()
				if err != nil {
					return err
				}
				continue
			}
			age := now.Sub(voteDate).Seconds()
			vote, _ := binary.Varint(v)
			expired := ((week - age) / week) * float64(vote)
			if math.Abs(expired) < 0.0001 {
				err = c.Delete()
				if err != nil {
					return err
				}
				continue
			}
			votes += expired
		}
		return nil
	})
	return votes, err
}

func getvotes(word string) (votes float64, err error) {
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(rankBucket)
		wb := b.Bucket([]byte(word))
		if wb == nil {
			votes = 0
			return nil
		}
		c := wb.Cursor()
		week := (time.Hour * 24 * 7).Seconds()
		now := time.Now()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			voteDate, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return err
			}
			if voteDate.Before(now.Add(time.Hour * -24 * 7)) {
				err = c.Delete()
				if err != nil {
					return err
				}
				continue
			}
			age := now.Sub(voteDate).Seconds()
			vote, _ := binary.Varint(v)
			expired := ((week - age) / week) * float64(vote)
			if math.Abs(expired) < 0.0001 {
				err = c.Delete()
				if err != nil {
					return err
				}
				continue
			}
			votes += expired
		}
		if wb.Stats().KeyN == 0 {
			err := b.DeleteBucket([]byte(word))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return votes, err
}

type ranking struct {
	name  string
	votes float64
}

type rankings []ranking

func getRanks() (ranks []ranking, err error) {
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(rankBucket)
		c := b.Cursor()
		week := (time.Hour * 24 * 7).Seconds()
		now := time.Now()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			votes := float64(0)
			c := b.Bucket(k).Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				voteDate, err := time.Parse(time.RFC3339, string(k))
				if err != nil {
					return err
				}
				if voteDate.Before(now.Add(time.Hour * -24 * 7)) {
					err = c.Delete()
					if err != nil {
						return err
					}
					continue
				}
				age := now.Sub(voteDate).Seconds()
				vote, _ := binary.Varint(v)
				expired := ((week - age) / week) * float64(vote)
				if math.Abs(expired) < 0.0001 {
					err = c.Delete()
					if err != nil {
						return err
					}
					continue
				}
				votes += expired
			}
			if b.Bucket(k).Stats().KeyN == 0 {
				err := b.DeleteBucket(k)
				if err != nil {
					return err
				}
			} else {
				ranks = append(ranks, ranking{string(k), votes})
			}
		}
		return nil
	})
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].votes > ranks[j].votes
	})
	return ranks, err
}

func addNickHostmask(hostmask, nick string) error {
	nick = strings.TrimRight(nick, "_")
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(hostmaskBucket)
		hb, err := b.CreateBucketIfNotExists([]byte(hostmask))
		if err != nil {
			return err
		}
		c := hb.Cursor()
		key, value := c.Last()
		if key != nil && string(value) == nick {
			return nil
		}
		return hb.Put([]byte(time.Now().UTC().Format(time.RFC3339)), []byte(nick))
	})
	return err
}

func addQuiet(ip string, dur time.Duration) (err error) {
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(quietBucket)
		if v := b.Get([]byte(ip)); v != nil {
			return nil
		}
		return b.Put([]byte(ip), []byte(time.Now().Add(dur).UTC().Format(time.RFC3339)))
	})
	return err
}

func getQuiet(ip string) (t time.Time, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(quietBucket)
		if v := b.Get([]byte(ip)); v != nil {
			t, err = time.Parse(time.RFC3339, string(v))
			return err
		}
		return fmt.Errorf("no entry found")
	})
	return
}

func removeQuiet(ip string) (err error) {
	err = db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket(quietBucket)
		return b.Delete([]byte(ip))
	})
	return
}
