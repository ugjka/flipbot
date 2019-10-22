package main

import (
	"fmt"
	"strings"

	wolf "github.com/Krognol/go-wolfram"
	hbot "github.com/ugjka/hellabot"
)

var calc = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!calc ")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		w := &wolf.Client{AppID: wolframAPIKey}
		res, err := w.GetShortAnswerQuery(m.Content[6:], wolf.Metric, 10)
		if err != nil {
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return true
		}
		irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, res))
		return false
	},
}
