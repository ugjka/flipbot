package main

import (
	"regexp"
	"time"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var isdeadReg = regexp.MustCompile(`(?i).*!+(?:(?:is)?dead+.*).*`)
var isdead = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && isdeadReg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		d, err := getDead(
			time.Minute*15,
			time.Minute*30,
			time.Hour*1,
			time.Hour*2,
		)
		if err != nil {
			log.Error("getDead", "error", err)
			return false
		}
		irc.Reply(m, d.String())
		return false
	},
}
