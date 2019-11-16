package main

import (
	"encoding/json"
	"fmt"
	"strings"
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
