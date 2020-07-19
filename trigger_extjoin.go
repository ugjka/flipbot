package main

import (
	"strings"
	"sync"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var extjoinOnce = &sync.Once{}

var extJoin = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "JOIN" || m.Command == "CAP"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		if m.Command == "JOIN" {
			extjoinOnce.Do(func() {
				log.Info("cap", "extended-join account-notify", "requesting")
				irc.Send("CAP REQ :extended-join account-notify")
			})
		}
		if m.Command == "CAP" && strings.TrimSpace(m.Content) == "extended-join account-notify" && len(m.Params) > 1 && m.Params[1] == "ACK" {
			log.Info("cap", "extended-join account-notify", "got ack")
			irc.Send("CAP END")
			extJoinEnabled = true
		}
		return false
	},
}
