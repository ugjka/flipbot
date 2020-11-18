package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	kitty "flipbot/kittybot"
)

var tailTrig = regexp.MustCompile("(?i)^\\s*!+tail\\w*\\s+(?:(\\d+)\\s+)?([A-Za-z_\\-\\[\\]\\^{}|`][A-Za-z0-9_\\-\\[\\]\\^{}|`]{0,15}\\*?)$")

const maxTail = 10

var tail = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.To == ircChannel && m.Command == "PRIVMSG" && tailTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		capped := false
		tailLen := 5
		nick := ""
		matches := tailTrig.FindStringSubmatch(m.Content)
		if len(matches) == 3 {
			tmp, err := strconv.Atoi(matches[1])
			if err == nil {
				if tmp <= maxTail {
					tailLen = tmp
				} else {
					capped = true
					tailLen = maxTail
				}
			}
			nick = matches[2]
		} else {
			nick = matches[1]
		}
		nick = strings.ToLower(nick)
		var err error
		if strings.HasSuffix(nick, "*") {
			nick, _, err = getSeenPrefix(strings.TrimRight(nick, "*"))
			switch {
			case err == errNotSeen:
				bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
				return
			case err != nil:
				bot.Crit("tailTrig", "error", err)
				return
			}
		}
		msgs, err := userTail(nick, "!tail", tailLen)
		switch {
		case err == errNotSeen || err == errNoResults:
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return
		case err != nil:
			bot.Crit("tailTrig", "error", err)
			return
		}
		if capped {
			bot.Reply(m, fmt.Sprintf("%s: tail lenght is capped to %d", m.Name, maxTail))
		}
		for _, msg := range msgs {
			bot.Reply(m, fmt.Sprintf("[%s] [%s] %s\n", msg.Time.Format("2006-01-02 3:04PM MST"), msg.Nick, msg.Message))
		}
	},
}
