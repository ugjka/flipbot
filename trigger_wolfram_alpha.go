package main

import (
	"fmt"
	"regexp"

	wolf "github.com/Krognol/go-wolfram"
	kitty "github.com/ugjka/kittybot"
)

var calcTrig = regexp.MustCompile(`(?i)^\s*!+[ck]al[ck]\w*\s+(\S.*)$`)
var calc = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && calcTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		query := calcTrig.FindStringSubmatch(m.Content)[1]
		w := &wolf.Client{AppID: wolframAPIKey}

		res, err := w.GetShortAnswerQuery(query, wolf.Metric, 10)
		if err != nil {
			bot.Warn("calc", "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		bot.Reply(m, fmt.Sprintf("%s: %s", m.Name, limit(res, 1024)))
	},
}
