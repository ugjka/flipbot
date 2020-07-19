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
		extjoinOnce.Do(func() {
			log.Info("cap", "extended-join", "requesting")
			irc.Send("CAP REQ :extended-join")
		})
		if m.Command == "CAP" && strings.TrimSpace(m.Content) == "extended-join" && len(m.Params) > 1 && m.Params[1] == "ACK" {
			log.Info("cap", "extended-join", "got ack")
			irc.Send("CAP END")
		}
		return false
	},
}
