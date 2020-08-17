package main

import (
	"fmt"
	"sync"
	"time"

	kitty "github.com/ugjka/kittybot"
)

var nickickerTrig = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return (m.Command == "JOIN" || m.Command == "NICK") && m.Name != ircNick
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		hostmask := m.Prefix.User + "@" + m.Prefix.Host
		if m.Command == "JOIN" {
			err := addNickHostmask(hostmask, m.Name)
			if err != nil {
				bot.Crit("addNickHostmask", "error", err)
			}
			quiet, err := checkNickHostmask(hostmask, m.To)
			if err != nil {
				bot.Crit("checkNickHostmask", "error", err)
			}
			if quiet {
				const timeOut = time.Minute * 10
				ip := m.Prefix.Host
				t, err := getQuiet(ip)
				if err != nil {
					bot.Info("adding quiet", "ip", ip)
					err := addQuiet(ip, timeOut)
					if err != nil {
						bot.Crit("could not add quiet to db", "error", err)
						return
					}
					bot.Send(fmt.Sprintf("MODE %s +q *!*@%s", ircChannel, ip))
					bot.Send(fmt.Sprintf("NOTICE %s :you can talk after %s", m.Name, timeOut))
					time.AfterFunc(timeOut, func() {
						bot.Info("quiet timeout", "ip", ip)
						bot.Send(fmt.Sprintf("MODE %s -q *!*@%s", ircChannel, ip))
						err := removeQuiet(ip)
						if err != nil {
							bot.Crit("can't remove quiet", "error", err)
						}
					})
					return
				}
				if time.Now().UTC().After(t) {
					bot.Info("timout from db", "ip", ip)
					bot.Send(fmt.Sprintf("MODE %s -q *!*@%s", ircChannel, ip))
					err := removeQuiet(ip)
					if err != nil {
						bot.Crit("can't remove quiet", "error", err)
						return
					}
				} else {
					bot.Send(fmt.Sprintf("NOTICE %s :you can talk after %s", m.Name, t.Sub(time.Now())))
				}
			}
		}
		if m.Command == "NICK" {
			err := addNickHostmask(hostmask, m.Name)
			if err != nil {
				bot.Crit("addNickHostmask", "error", err)
			}
			kick, err := checkNickHostmask(hostmask, m.To)
			if err != nil {
				bot.Crit("checkNickHostmask", "error", err)
			}
			if kick {
				bot.Info("too many nick changes", "kicking", m.To)
				bot.Send(fmt.Sprintf("REMOVE %s %s :Too many nick changes in the past 24 hours", ircChannel, m.To))
			}
		}
	},
}

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
