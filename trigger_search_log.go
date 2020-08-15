package main

import (
	"fmt"
	"regexp"
	"strings"

	kitty "github.com/ugjka/kittybot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var searchLogRegex = regexp.MustCompile(`(?i)^\s*!search\s+(\S.*)$`)
var searchLog = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.To == ircChannel && m.Command == "PRIVMSG" && searchLogRegex.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		query := searchLogRegex.FindStringSubmatch(m.Content)[1]
		query = strings.ToLower(query)
		msgs, err := search(query, "!search")
		switch {
		case err == errNoResults:
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return
		case err != nil:
			log.Crit("searchLog", "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		for _, msg := range msgs {
			bot.Reply(m, fmt.Sprintf("[%s] [%s] %s", msg.Time.Format("2006-01-02 3:04PM MST"), msg.Nick, msg.Message))
		}
	},
}
