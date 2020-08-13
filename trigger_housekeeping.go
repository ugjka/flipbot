package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	kitty "github.com/ugjka/kittybot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var namesCall sync.Once

var names = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.To == ircChannel
	},
	Action: func(irc *kitty.Bot, m *kitty.Message) {
		namesCall.Do(func() {
			log.Info("firstrun", "action", "getting names")
			irc.Send("NAMES " + ircChannel)
		})
	},
}

var modes sync.Once
var setmodes = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.To == ircChannel
	},
	Action: func(irc *kitty.Bot, m *kitty.Message) {
		modes.Do(func() {
			go func(irc *kitty.Bot) {
				for {
					time.Sleep(time.Second * 30)
					irc.Send("PING " + ircServer)
				}
			}(irc)
			log.Info("setting modes for self", "modes", "+RQi")
			irc.Send(fmt.Sprintf("MODE %s +RQi", ircNick))
			time.AfterFunc(time.Second*60, func() {
				log.Info("failover", "action", "sending pass")
				irc.Msg("NickServ", fmt.Sprintf("IDENTIFY %s %s", ircNick, ircPassword))
			})
		})
	},
}

var voice = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "JOIN" || (len(m.Params) == 3 && m.Name == "ChanServ" && m.Command == "MODE" && (m.Params[1] == ("-v") || m.Params[1] == ("-o")))
	},
	Action: func(irc *kitty.Bot, m *kitty.Message) {
		if len(m.Params) == 3 && m.Name == "ChanServ" && m.Command == "MODE" {
			log.Info("giving voice", "to", m.Params[2], "in", m.To)
			irc.ChMode(m.Params[2], m.To, "+v")
			return
		}
		// hostmask := m.Prefix.User + "@" + m.Prefix.Host
		// quiet, err := checkNickHostmask(hostmask, m.To)
		// if err != nil {
		// 	log.Crit("checkNickHostmask", "error", err)
		// }
		// if quiet {
		// 	log.Info("too many nick changes, not voicing", "nick", m.Name, "hostmask", hostmask)
		// 	return false
		// }
		// ip := m.Prefix.Host
		// if _, err := getQuiet(ip); err == nil {
		// 	log.Info("nick quieted, not voicing", "nick", m.Name, "ip", ip)
		// 	return false
		// }
		log.Info("giving voice", "to", m.Name, "in", m.To)
		irc.ChMode(m.Name, m.To, "+v")
	},
}

var voicenames = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "353"
	},
	Action: func(irc *kitty.Bot, m *kitty.Message) {
		time.Sleep(time.Second * 5)
		for _, k := range strings.Split(m.Content, " ") {
			if strings.HasPrefix(k, "+") || strings.HasPrefix(k, "@") {
				continue
			}
			if irc.Nick == k {
				continue
			}
			log.Info("giving voice", "to", k, "in", m.Params[2])
			irc.ChMode(k, m.Params[2], "+v")
		}
	},
}
