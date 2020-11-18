package main

import (
	"sync"
	"time"

	kitty "flipbot/kittybot"
)

const nickChangeWindow = time.Hour * 24
const nickChangesMax = 6

var nickickerCleanupOnce = &sync.Once{}

var nickickerCleanupTrig = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PING" || m.Command == "PONG" || (m.Command == "PRIVMSG" && m.To == ircChannel)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		nickickerCleanupOnce.Do(func() {
			bot.Info("info", "starting quiet timers", "started")
			err := quietTimers(bot)
			if err != nil {
				bot.Crit("couln't start quiet timers", "error", err)
			}
			bot.Info("info", "starting quiet timers", "executed")
		})
	},
}
