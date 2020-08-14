package main

import (
	"strings"
	"sync"

	kitty "github.com/ugjka/kittybot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var extjoinOnce = &sync.Once{}

var extJoin = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "JOIN" || m.Command == "CAP"
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		if m.Command == "JOIN" {
			extjoinOnce.Do(func() {
				log.Info("cap", "extended-join account-notify", "requesting")
				b.Send("CAP REQ :extended-join account-notify")
			})
		}
		if m.Command == "CAP" && strings.TrimSpace(m.Content) == "extended-join account-notify" && len(m.Params) > 1 && m.Params[1] == "ACK" {
			log.Info("cap", "extended-join account-notify", "got ack")
			b.Send("CAP END")
			extJoinEnabled = true
		}
	},
}
