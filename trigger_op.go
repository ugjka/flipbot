package main

import (
	"strings"
	"time"

	hbot "github.com/ugjka/hellabot"
	"github.com/ugjka/reverse"
	gomail "gopkg.in/gomail.v2"
	log "gopkg.in/inconshreveable/log15.v2"
)

var notifyop = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel &&
			strings.Contains(m.Content, op)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		history := ""
		logCTR.Lock()
		scan := reverse.NewScanner(logCTR.File)
		for i := 0; i < 20; i++ {
			if scan.Scan() {
				history = scan.Text() + "\n" + history
			}
		}
		logCTR.Unlock()
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
