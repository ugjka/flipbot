package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var tailTrig = regexp.MustCompile("(?i)^\\s*!+tail\\w*\\s+(?:(\\d+)\\s+)?([A-Za-z_\\-\\[\\]\\^{}|`][A-Za-z0-9_\\-\\[\\]\\^{}|`]{0,15}\\*?)$")

const maxTail = 15

var tail = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && tailTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
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
				irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
				return false
			case err != nil:
				log.Crit("tailTrig", "error", err)
				return false
			}
		}
		msgs, err := userTail(nick, "!tail", tailLen)
		switch {
		case err == errNotSeen || err == errNoResults:
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return false
		case err != nil:
			log.Crit("tailTrig", "error", err)
			return false
		}
		if capped {
			irc.Reply(m, fmt.Sprintf("%s: tail lenght is capped to %d", m.Name, maxTail))
		}
		for _, msg := range msgs {
			irc.Reply(m, fmt.Sprintf("[%s] [%s] %s\n", msg.Time.Format("2006-01-02 3:04PM MST"), msg.Nick, msg.Message))
		}
		return false
	},
}
