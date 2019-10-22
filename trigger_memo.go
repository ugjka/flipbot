package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

type memoStruct struct {
	Sender  string
	Message string
}

var memoCTR = struct {
	store map[string][]memoStruct
	*os.File
	sync.RWMutex
}{
	store: make(map[string][]memoStruct),
}

var memo = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!memo ")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		msg := strings.TrimPrefix(m.Content, "!memo ")
		msg = whitespace.ReplaceAllString(msg, " ")
		args := strings.Split(msg, " ")
		nick := strings.ToLower(args[0])
		msg = strings.Join(args[1:], " ")
		memoCTR.Lock()
		defer memoCTR.Unlock()
		memoCTR.store[nick] = append(memoCTR.store[nick], memoStruct{m.Name, msg})
		tmp, err := json.Marshal(memoCTR.store)
		if err == nil {
			err := memoCTR.Truncate(0)
			if err != nil {
				log.Crit("Could not truncate memo file", "error", err)
				return false
			}
			if _, err = memoCTR.WriteAt(tmp, 0); err == nil {
				irc.Reply(m, fmt.Sprintf("%s's memo will be sent to %s when I see them join or post", m.Name, nick))
				return false
			}
			log.Crit("Could not write to memo file in memo", "error", err)

		}
		irc.Reply(m, fmt.Sprintf("error: %v", err))
		return false
	},
}

var memowatcher = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "JOIN" || m.Command == "PRIVMSG"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		memoCTR.Lock()
		defer memoCTR.Unlock()
		if v, ok := memoCTR.store[strings.ToLower(m.Name)]; ok {
			for _, v := range v {
				irc.Msg(ircChannel, fmt.Sprintf("%s's memo to %s: %s", v.Sender, m.Name, v.Message))
			}
			delete(memoCTR.store, strings.ToLower(m.Name))
			tmp, err := json.Marshal(memoCTR.store)
			if err == nil {
				err := memoCTR.Truncate(0)
				if err != nil {
					log.Crit("Could not truncate memowatcher", "error", err)
					return false
				}
				if _, err := memoCTR.WriteAt(tmp, 0); err != nil {
					log.Crit("Could not write to memo file", "error", err)
					return false
				}
			}
		}
		return false
	},
}
