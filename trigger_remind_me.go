package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hako/durafmt"

	hbot "github.com/ugjka/hellabot"
	"github.com/ugjka/remindme"
	log "gopkg.in/inconshreveable/log15.v2"
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
						onlineCTR.RLock()
						if _, ok := onlineCTR.db[x.Name]; !ok {
							memoCTR.Lock()
							memoCTR.store[strings.ToLower(x.Name)] = append(memoCTR.store[strings.ToLower(x.Name)],
								memoStruct{ircNick, fmt.Sprintf("Your reminder: %s", x.Message)})
							tmp, err := json.Marshal(memoCTR.store)
							if err == nil {
								err := memoCTR.Truncate(0)
								if err != nil {
									log.Crit("Could not truncate memo file", "error", err)
									goto exit
								}
								if _, err := memoCTR.WriteAt(tmp, 0); err == nil {
									goto exit
								} else {
									log.Crit("Could not write to memo file in memo", "error", err)
								}
							exit:
							}
							memoCTR.Unlock()
						}
						onlineCTR.RUnlock()
					}
				}
			}()
		})
		return false
	},
}

var getreminder = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.To == ircChannel && m.Command == "PRIVMSG" &&
			strings.HasPrefix(m.Content, "!remindme ")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		r, err := remindme.Parse(m.Content[10:])
		if err != nil {
			irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, err))
			return false
		}
		r.Name = m.Name
		remind.Add(r)
		delta := r.Target.Sub(time.Now())
		dur := durafmt.Parse(delta)
		irc.Reply(m, fmt.Sprintf("%s: Your reminder will fire %s from now", m.Name, removeMilliseconds(dur.String())))
		return false
	},
}
