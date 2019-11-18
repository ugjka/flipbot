package main

import (
	"fmt"
	"regexp"
	"strings"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var searchLogRegex = regexp.MustCompile(`(?i)^\s*!search\s+(\S.*)$`)
var searchLog = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && searchLogRegex.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, fmt.Sprintf("%s: this feature is disabled due to negative feedback", m.Name))
		return false
		query := searchLogRegex.FindStringSubmatch(m.Content)[1]
		query = strings.ToLower(query)
		msgs, err := search(query, "!search")
		switch {
		case err == errNoResults:
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return false
		case err != nil:
			log.Crit("searchLog", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return false
		}
		for _, msg := range msgs {
			irc.Reply(m, fmt.Sprintf("[%s] [%s] %s", msg.Time.Format("2006-01-02 3:04PM MST"), msg.Nick, msg.Message))
		}
		return false
	},
}
