package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	hbot "github.com/ugjka/hellabot"
	gomail "gopkg.in/gomail.v2"
	log "gopkg.in/inconshreveable/log15.v2"
)

var notifyop = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel &&
			strings.Contains(m.Content, op) && strings.ToLower(m.Name) != op
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		riga, err := time.LoadLocation("Europe/Riga")
		if err != nil {
			log.Crit("notifyop", "error", err)
			return false
		}
		history := ""
		err = db.View(func(tx *bolt.Tx) error {
			c := tx.Bucket(logBucket).Cursor()
			i := 0
			msg := &Message{}
			for k, v := c.Last(); k != nil && v != nil; k, v = c.Prev() {
				if i > 20 {
					break
				}
				i++
				err := json.Unmarshal(v, &msg)
				if err != nil {
					return err
				}
				history = fmt.Sprintf("%s <%s> %s\n", msg.Time.In(riga).Format(time.Kitchen), msg.Nick, msg.Message) + history
			}
			return nil
		})
		if err != nil {
			log.Crit("notifyop", "error", err)
			return false
		}
		msg := gomail.NewMessage()
		msg.SetHeader("From", serverEmail)
		msg.SetHeader("To", email)
		msg.SetHeader("Subject", "irc notification from "+m.Name)
		msg.SetBody("text/plain", "--------------------------\n"+
			m.Name+": "+m.Content+"\n"+
			"--------------------------\n"+
			time.Now().String()+"\n\n"+
			"HISTORY:\n"+
			history)

		d := gomail.NewDialer("127.0.0.1", 25, "", "")

		if err := d.DialAndSend(msg); err != nil {
			log.Crit("could not push op nick highlight", "error", err)
			return false
		}
		return false
	},
}

var logOnce = &sync.Once{}
var indexLog = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.From == op && m.Content == "!index"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		logOnce.Do(func() {
			var max uint64
			err := db.View(func(tx *bolt.Tx) error {
				last, _ := tx.Bucket(logBucket).Cursor().Last()
				max = btoi(last)
				return nil
			})
			if err != nil {
				log.Crit("indexLog", "error", err)
				return
			}
			semaphore := make(chan struct{}, 100)
			for i := 1; i < int(max); i++ {
				if i%10000 == 0 {
					irc.Msg(op, fmt.Sprintf("indexing %d out of %d", i, max))
				}
				semaphore <- struct{}{}
				go func(i int) {
					err := db.Batch(func(tx *bolt.Tx) error {
						index := tx.Bucket(indexBucket)
						v := tx.Bucket(logBucket).Get(itob(uint64(i)))
						if v == nil {
							return fmt.Errorf("nil value")
						}
						msg := Message{}
						err := json.Unmarshal(v, &msg)
						if err != nil {
							return err
						}
						for _, v := range split(strings.ToLower(msg.Message)) {
							b, err := index.CreateBucketIfNotExists([]byte(v))
							if err != nil {
								return err
							}
							err = b.Put(itob(uint64(i)), []byte(""))
							if err != nil {
								return err
							}
						}
						return err
					})
					if err != nil {
						log.Crit("indexLog", "error", err)
					}
					<-semaphore
				}(i)
			}
			irc.Msg(op, "Indexing done!!!! Hip Hip Hurray!!!")
			return
		})
		return false
	},
}

var usersOnce = &sync.Once{}
var indexUsers = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.From == op && m.Content == "!indexusers"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		usersOnce.Do(func() {
			var max uint64
			err := db.View(func(tx *bolt.Tx) error {
				last, _ := tx.Bucket(logBucket).Cursor().Last()
				max = btoi(last)
				return nil
			})
			if err != nil {
				log.Crit("indexUsers", "error", err)
				return
			}
			semaphore := make(chan struct{}, 100)
			for i := 1; i < int(max); i++ {
				if i%10000 == 0 {
					irc.Msg(op, fmt.Sprintf("\r%d out of %d", i, max))
				}
				semaphore <- struct{}{}
				go func(i int) {
					err := db.Batch(func(tx *bolt.Tx) error {
						users := tx.Bucket(usersBucket)
						v := tx.Bucket(logBucket).Get(itob(uint64(i)))
						if v == nil {
							return fmt.Errorf("nil value")
						}
						msg := Message{}
						err := json.Unmarshal(v, &msg)
						if err != nil {
							return err
						}
						b, err := users.CreateBucketIfNotExists([]byte(strings.ToLower(msg.Nick)))
						if err != nil {
							return err
						}
						err = b.Put(itob(uint64(i)), []byte(""))
						if err != nil {
							return err
						}

						return err
					})
					if err != nil {
						log.Crit("indexUsers", "error", err)
					}
					<-semaphore
				}(i)
			}
			irc.Msg(op, "Users Indexed!!! Hip Hip Hurray!!!")
		})
		return false
	},
}
