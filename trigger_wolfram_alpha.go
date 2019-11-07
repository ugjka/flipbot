package main

import (
	"fmt"
	"regexp"

	wolf "github.com/Krognol/go-wolfram"
	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var calcTrig = regexp.MustCompile(`(?i)^\s*!+[ck]al[ck]\w*\s+(\S.*)$`)
var calc = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && calcTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		query := calcTrig.FindStringSubmatch(m.Content)[1]
		w := &wolf.Client{AppID: wolframAPIKey}

		res, err := w.GetShortAnswerQuery(query, wolf.Metric, 10)
		if err != nil {
			log.Warn("calc", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return true
		}
		irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, limit(res)))
		return false
	},
}
