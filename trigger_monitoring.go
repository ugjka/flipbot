package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	"github.com/hako/durafmt"
	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var logCTR = struct {
	*os.File
	sync.Mutex
}{}

//Seen struct
type Seen struct {
	Seen    time.Time
	LastMSG string
}

var seenTrig = regexp.MustCompile("(?i)^\\s*!+seen\\w*\\s+([A-Za-z_\\-\\[\\]\\^{}|`][A-Za-z0-9_\\-\\[\\]\\^{}|`]{0,15}\\*?)$")

var seen = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && seenTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		nick := seenTrig.FindStringSubmatch(m.Content)[1]
		nick = strings.ToLower(nick)
		if nick == strings.ToLower(m.Name) {
			irc.Reply(m, fmt.Sprintf("%s: I'm seeing you!", m.Name))
			return false
		}
		seen := Seen{}
		var err error
		if strings.HasSuffix(nick, "*") {
			nick, seen, err = getSeenPrefix(strings.TrimRight(nick, "*"))
		} else {
			seen, err = getSeen(nick)
		}
		switch {
		case err == errNotSeen:
			irc.Reply(m, fmt.Sprintf("%s: I haven't seen that nick before", m.Name))
			return false
		case err != nil:
			log.Warn("getSeen", "error", err)
			return false
		}
		dur := durafmt.Parse(time.Now().UTC().Sub(seen.Seen))
		if seen.LastMSG != "" {
			irc.Reply(m, fmt.Sprintf("%s: I saw %s %s ago. Their last message was: %s", m.Name, nick, roundDuration(dur.String()), seen.LastMSG))
		} else {
			irc.Reply(m, fmt.Sprintf("%s: I saw %s %s ago", m.Name, nick, roundDuration(dur.String())))
		}

		return false
	},
}

func roundDuration(dur string) string {
	arr := strings.Split(dur, " ")
	if len(arr) > 2 {
		return strings.Join(arr[:4], " ")
	}
	return strings.Join(arr[:2], " ")
}

var watcher = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return (m.Command == "PRIVMSG" && m.To == ircChannel) || m.Command == "JOIN" ||
			m.Command == "QUIT" || m.Command == "PART" || m.Command == "KICK"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		name := ""
		if m.Command == "KICK" {
			name = m.Params[1]
		} else {
			name = m.Name
		}
		name = strings.ToLower(name)
		seen, err := getSeen(name)
		switch {
		case err == errNotSeen:
			break
		case err != nil:
			log.Warn("getSeen", "error", err)
			return false
		}
		if m.Command == "PRIVMSG" {
			seen.Seen = time.Now().UTC()
			seen.LastMSG = m.Content
			err := setSeen(name, &seen)
			if err != nil {
				log.Warn("setSeen", "error", err)
				return false
			}
		} else {
			seen.Seen = time.Now().UTC()
			err := setSeen(name, &seen)
			if err != nil {
				log.Warn("setSeen", "error", err)
				return false
			}
		}
		return false
	},
}

var topTrig = regexp.MustCompile(`(?i).*!+(?:top|masters?|masterminds?|echoline).*`)
var top = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && topTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		stats := make(map[string]int)
		res := make(result, 0)
		week := time.Now().UTC().Add(time.Hour * -24 * 7)
		msg := Message{}
		total := 0
		err := db.View(func(tx *bolt.Tx) error {
			c := tx.Bucket(logBucket).Cursor()
			for k, v := c.Last(); k != nil && v != nil; k, v = c.Prev() {
				err := json.Unmarshal(v, &msg)
				if err != nil {
					return err
				}
				if msg.Time.Before(week) {
					break
				}
				total++
				stats[strings.ToLower(msg.Nick)]++
			}
			return nil
		})
		if err != nil {
			log.Crit("!top", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return false
		}
		for k, v := range stats {
			k = k[:len(k)-1] + "*"
			res = append(res, stat{k, v})
		}
		sort.Sort(sort.Reverse(res))
		out := "Top posters for past 7 days: "
		out = fmt.Sprintf("Daily Average: %d posts. %s", total/7, out)
		for i, v := range res {
			if (i == 9) || (i == len(res)-1) {
				out += fmt.Sprintf("%d. %s posts.", i+1, v)
				break
			}
			out += fmt.Sprintf("%d. %s posts, ", i+1, v)
		}
		irc.Reply(m, out)
		return false
	},
}

var logmsg = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		logCTR.Lock()
		fmt.Fprintf(logCTR.File, "[%s] <%s>\t%s\n", time.Now().UTC().Format("06:01:02|15:04:05"), m.Name, m.Content)
		logCTR.Unlock()
		return false
	},
}

var logJoin = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "JOIN"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		logCTR.Lock()
		account := "[unknown]"
		if len(m.Params) == 3 {
			account = m.Params[1]
		}
		fmt.Fprintf(logCTR.File, "[%s] [JOIN]\t%s!%s@%s (%s) account: %s\n", time.Now().UTC().Format("06:01:02|15:04:05"), m.Name, m.User, m.Host, m.Trailing(), account)
		logCTR.Unlock()
		return false
	},
}

var logAccount = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "ACCOUNT"
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		logCTR.Lock()
		fmt.Fprintf(logCTR.File, "[%s] [ACCOUNT]\t%s!%s@%s (%s)\n", time.Now().UTC().Format("06:01:02|15:04:05"), m.Name, m.User, m.Host, m.Params[0])
		logCTR.Unlock()
		return false
	},
}

var logmsgBolt = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		err := setLogMSG(&Message{
			Time:    time.Now(),
			Nick:    m.Name,
			Message: m.Content,
		})
		if err != nil {
			log.Crit("setLogMSG", "error", err)
		}
		return false
	},
}

//Message is an irc message
type Message struct {
	Time    time.Time
	Nick    string
	Message string
}

type stat struct {
	name  string
	count int
}

type result []stat

func (r result) Len() int {
	return len(r)
}

func (r result) Less(i, j int) bool {
	if r[i].count == r[j].count {
		return !sort.StringsAreSorted([]string{strings.ToLower(r[i].name), strings.ToLower(r[j].name)})
	}
	if r[i].count < r[j].count {
		return true
	}
	return false
}

func (r result) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (s stat) String() string {
	return fmt.Sprintf("%s:\t%d\t", s.name, s.count)
}
func (r result) String() (out string) {
	for _, v := range r {
		out += fmt.Sprintf("%s\n", v)
	}
	return
}
