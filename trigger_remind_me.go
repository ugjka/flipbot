package main

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/hako/durafmt"

	hbot "github.com/ugjka/hellabot"
	"github.com/ugjka/remindme"
)

var remindOnce sync.Once
var reminder = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.To == ircChannel
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		remindOnce.Do(func() {
			go func() {
				for {
					select {
					case x := <-remind.Receive:
						text := fmt.Sprintf("%s's reminder: %s", x.Name, x.Message)
						irc.Msg(ircChannel, text)
					}
				}
			}()
		})
		return false
	},
}

var getreminderTrig = regexp.MustCompile(`(?i)^\s*!+remind(?:er)?\s+(\S.*)$`)
var getreminder = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.To == ircChannel && m.Command == "PRIVMSG" &&
			getreminderTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		r, err := remindme.Parse(getreminderTrig.FindStringSubmatch(m.Content)[1])
		if err != nil {
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return false
		}
		r.Name = m.Name
		remind.Add(r)
		delta := r.Target.Sub(time.Now())
		dur := durafmt.Parse(delta)
		irc.Reply(m, fmt.Sprintf("%s: Your reminder will fire %s from now", m.Name, roundDuration(dur.String())))
		return false
	},
}
