package main

import (
	"fmt"
	"regexp"
	"time"

	kitty "bootybot/kittybot"

	gomail "gopkg.in/gomail.v2"
)

var notifyopReg = regexp.MustCompile(`(?i).*!+(?:op+|alarm+|alert+).*`)
var notifyop = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && notifyopReg.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		msg := gomail.NewMessage()
		msg.SetHeader("From", msg.FormatAddress(serverEmail, "mbd"))
		msg.SetHeader("To", email)
		msg.SetHeader("Subject", "irc notification from "+m.Name)
		msg.SetBody("text/plain", "--------------------------\n"+
			m.Name+": "+m.Content+"\n\n"+
			fmt.Sprintf("%s\n", m.Name)+
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
