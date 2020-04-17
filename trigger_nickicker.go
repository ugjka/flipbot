package main

import (
	"fmt"
	"time"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var nickickerTrig = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return (m.Command == "JOIN" || m.Command == "NICK") && m.Name != ircNick
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		hostmask := m.Prefix.User + "@" + m.Prefix.Host
		if m.Command == "JOIN" {
			err := addNickHostmask(hostmask, m.Name)
			if err != nil {
				log.Crit("addNickHostmask", "error", err)
			}
		}
		if m.Command == "NICK" {
			err := addNickHostmask(hostmask, m.Name)
			if err != nil {
				log.Crit("addNickHostmask", "error", err)
			}
			kick, err := checkNickHostmask(hostmask, m.To)
			if err != nil {
				log.Crit("checkNickHostmask", "error", err)
			}
			if kick {
				log.Info("too many nick changes", "kicking", m.To)
				irc.Send(fmt.Sprintf("REMOVE %s %s :Too many nick changes in the past 24 hours", ircChannel, m.To))
			}
		}
		return false
	},
}

const nickChangeWindow = time.Hour * 24
const nickChangesMax = 6
