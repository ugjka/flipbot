package main

import (
	"regexp"
	"time"

	kitty "github.com/ugjka/kittybot"
)

var isDead = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		var isDeadReg = regexp.MustCompile(`(?i).*!+(?:(?:is)?dead+.*).*`)
		return m.Command == "PRIVMSG" && isDeadReg.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		d, err := getDead(
			time.Minute*15,
			time.Minute*30,
			time.Hour*1,
			time.Hour*2,
		)
		if err != nil {
			bot.Error("getDead", "error", err)
			return
		}
		bot.Reply(m, d.String())
	},
}

var isRecent = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		var isRecentReg = regexp.MustCompile(`(?i).*!+(?:(?:is)?recent+.*).*`)
		return m.Command == "PRIVMSG" && isRecentReg.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		r, err := getRecent(10)
		if err != nil {
			bot.Error("getRecent", "error", err)
			return
		}
		bot.Reply(m, r.String())
	},
}
