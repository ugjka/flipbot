package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var modes sync.Once
var setmodes = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.To == ircChannel
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		modes.Do(func() {
			go func(irc *hbot.Bot) {
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
		return false
	},
}

var voice = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "JOIN" || (len(m.Params) == 3 && m.Name == "ChanServ" && m.Command == "MODE" && (m.Params[1] == ("-v") || m.Params[1] == ("-o")))
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		if len(m.Params) == 3 && m.Name == "ChanServ" && m.Command == "MODE" {
			log.Info("giving voice", "to", m.Params[2], "in", m.To)
			irc.ChMode(m.Params[2], m.To, "+v")
			return false
		}
		log.Info("giving voice", "to", m.Name, "in", m.To)
		irc.ChMode(m.Name, m.To, "+v")
		return false
	},
}

var voicenames = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "353"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
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
		return false
	},
}
