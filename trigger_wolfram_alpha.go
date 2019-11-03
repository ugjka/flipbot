package main

import (
	"fmt"
	"strings"

	wolf "github.com/Krognol/go-wolfram"
	hbot "github.com/ugjka/hellabot"
)

const calcTrig = "!calc "

var calc = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, calcTrig)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		query := strings.TrimPrefix(m.Content, calcTrig)
		w := &wolf.Client{AppID: wolframAPIKey}

		res, err := w.GetShortAnswerQuery(query, wolf.Metric, 10)
		if err != nil {
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return true
		}
		irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, limit(res)))
		return false
	},
}
