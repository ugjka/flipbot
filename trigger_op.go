package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	kitty "github.com/ugjka/kittybot"
	gomail "gopkg.in/gomail.v2"
)

var notifyopReg = regexp.MustCompile(`(?i).*!+(?:op+|alarm+|alert+).*`)
var notifyop = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel && notifyopReg.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		msg := gomail.NewMessage()
		msg.SetHeader("From", msg.FormatAddress(serverEmail, "mbd"))
		msg.SetHeader("To", email)
		msg.SetHeader("Subject", "irc notification from "+m.Name)
		msg.SetBody("text/plain", "--------------------------\n"+
			m.Name+": "+m.Content+"\n\n"+
			fmt.Sprintf("%s!%s@%s\n", m.Name, m.User, m.Host)+
			"--------------------------\n"+
			time.Now().String()+"\n\n")

		d := gomail.NewDialer("127.0.0.1", 25, "", "")

		if err := d.DialAndSend(msg); err != nil {
			bot.Crit("could not push op nick highlight", "error", err)
			return
		}
		bot.Reply(m, fmt.Sprintf("%s: message emailed to %s", m.Name, op))
	},
}

var logOnce = &sync.Once{}
var indexLog = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.From == op && m.Content == "!index"
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		logOnce.Do(func() {
			var max uint64
			err := db.View(func(tx *bolt.Tx) error {
				last, _ := tx.Bucket(logBucket).Cursor().Last()
				max = btoi(last)
				return nil
			})
			if err != nil {
				bot.Crit("indexLog", "error", err)
				return
			}
			semaphore := make(chan struct{}, 100)
			for i := 1; i < int(max); i++ {
				if i%10000 == 0 {
					bot.Msg(op, fmt.Sprintf("indexing %d out of %d", i, max))
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
						bot.Crit("indexLog", "error", err)
					}
					<-semaphore
				}(i)
			}
			bot.Msg(op, "Indexing done!!!! Hip Hip Hurray!!!")
			return
		})
	},
}

var usersOnce = &sync.Once{}
var indexUsers = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.From == op && m.Content == "!indexusers"
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		usersOnce.Do(func() {
			var max uint64
			err := db.View(func(tx *bolt.Tx) error {
				last, _ := tx.Bucket(logBucket).Cursor().Last()
				max = btoi(last)
				return nil
			})
			if err != nil {
				bot.Crit("indexUsers", "error", err)
				return
			}
			semaphore := make(chan struct{}, 100)
			for i := 1; i < int(max); i++ {
				if i%10000 == 0 {
					bot.Msg(op, fmt.Sprintf("\r%d out of %d", i, max))
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
						bot.Crit("indexUsers", "error", err)
					}
					<-semaphore
				}(i)
			}
			bot.Msg(op, "Users Indexed!!! Hip Hip Hurray!!!")
		})
	},
}
