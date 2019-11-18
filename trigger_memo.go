package main

import (
	"fmt"
	"regexp"
	"strings"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

type memoItem struct {
	Sender  string
	Message string
}

type memos []memoItem

var memoTrig = regexp.MustCompile("(?i)^\\s*!+memo\\w*\\s+([A-Za-z_\\-\\[\\]\\^{}|`][A-Za-z0-9_\\-\\[\\]\\^{}|`]{1,15})\\s+(\\S.+)$")
var memo = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.To == ircChannel && m.Command == "PRIVMSG" && memoTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		matches := memoTrig.FindStringSubmatch(m.Content)
		nick := strings.ToLower(matches[1])
		msg := matches[2]
		err := setMemo(nick, memoItem{m.Name, msg})
		if err != nil {
			log.Crit("setMemo", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s's memo will be sent to %s when I see them join or post", m.Name, nick))
		return false
	},
}

var memowatcher = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "JOIN" || m.Command == "PRIVMSG"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		memos, err := getMemo(strings.ToLower(m.Name))
		switch {
		case err == errNoMemo:
			return false
		case err != nil:
			log.Crit("getMemo", "error", err)
			return false
		}
		for _, v := range memos {
			irc.Msg(ircChannel, fmt.Sprintf("%s's memo to %s: %s", v.Sender, m.Name, v.Message))
		}
		return false
	},
}
