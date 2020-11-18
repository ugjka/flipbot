package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	kitty "flipbot/kittybot"
)

var setup = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "JOIN" && m.Name == ircNick
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		bot.Info("getting NAMES")
		bot.Send("NAMES " + ircChannel)
		bot.Info("setting user modes", "modes", "+RQi")
		bot.Send(fmt.Sprintf("MODE %s +RQi", ircNick))
	},
}

type pinger struct {
	once sync.Once
}

func (p *pinger) Handle(bot *kitty.Bot, m *kitty.Message) {
	p.once.Do(func() {
		go func(bot *kitty.Bot) {
			for {
				time.Sleep(time.Second * 30)
				bot.Send("PING " + ircServer)
			}
		}(bot)
	})
}

var voice = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "JOIN" || (len(m.Params) == 3 && m.Name == "ChanServ" && m.Command == "MODE" && (m.Params[1] == ("-v") || m.Params[1] == ("-o")))
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		if len(m.Params) == 3 && m.Name == "ChanServ" && m.Command == "MODE" {
			bot.Info("giving voice", "to", m.Params[2], "in", m.To)
			bot.ChMode(m.Params[2], m.To, "+v")
			return
		}
		// hostmask := m.Prefix.User + "@" + m.Prefix.Host
		// quiet, err := checkNickHostmask(hostmask, m.To)
		// if err != nil {
		// 	bot.Crit("checkNickHostmask", "error", err)
		// }
		// if quiet {
		// 	bot.Info("too many nick changes, not voicing", "nick", m.Name, "hostmask", hostmask)
		// 	return false
		// }
		// ip := m.Prefix.Host
		// if _, err := getQuiet(ip); err == nil {
		// 	bot.Info("nick quieted, not voicing", "nick", m.Name, "ip", ip)
		// 	return false
		// }
		bot.Info("giving voice", "to", m.Name, "in", m.To)
		bot.ChMode(m.Name, m.To, "+v")
	},
}

var voicenames = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "353"
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		time.Sleep(time.Second * 5)
		for _, k := range strings.Split(m.Content, " ") {
			if strings.HasPrefix(k, "+") || strings.HasPrefix(k, "@") {
				continue
			}
			if bot.Nick == k {
				continue
			}
			bot.Info("giving voice", "to", k, "in", m.Params[2])
			bot.ChMode(k, m.Params[2], "+v")
		}
	},
}
