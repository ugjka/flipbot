package main

import (
	"fmt"
	"strings"
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

const getreminderTrig = "!reminder "

var getreminder = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.To == ircChannel && m.Command == "PRIVMSG" &&
			strings.HasPrefix(m.Content, getreminderTrig)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		r, err := remindme.Parse(strings.TrimPrefix(m.Content, getreminderTrig))
		if err != nil {
			irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, err))
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
