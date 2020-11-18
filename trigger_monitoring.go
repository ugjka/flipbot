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

	kitty "flipbot/kittybot"

	"github.com/boltdb/bolt"
	"github.com/hako/durafmt"
)

var logCTR = struct {
	*os.File
	sync.Mutex
}{}

//Seen struct
type Seen struct {
	Seen    time.Time
	LastMSG string
	Command string
}

var seenTrig = regexp.MustCompile("(?i)^\\s*!+seen\\w*\\s+([A-Za-z_\\-\\[\\]\\^{}|`][A-Za-z0-9_\\-\\[\\]\\^{}|`]{0,15}\\*?)$")

var seen = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && seenTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		nick := seenTrig.FindStringSubmatch(m.Content)[1]
		nick = strings.ToLower(nick)
		if nick == strings.ToLower(m.Name) {
			bot.Reply(m, fmt.Sprintf("%s: I'm seeing you!", m.Name))
			return
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
			bot.Reply(m, fmt.Sprintf("%s: I haven't seen that nick before", m.Name))
			return
		case err != nil:
			bot.Warn("getSeen", "error", err)
			return
		}
		dur := durafmt.Parse(time.Now().UTC().Sub(seen.Seen))
		msg := fmt.Sprintf("%s: I saw %s %s ago. ", m.Name, nick, roundDuration(dur.String()))
		if seen.Command != "" {
			msg += fmt.Sprintf("Last activity: %s. ", seen.Command)
		}
		if seen.LastMSG != "" {
			msg += fmt.Sprintf("Last message: %s", seen.LastMSG)
		}
		msg = limitReply(bot, m, msg, 1)
		bot.Reply(m, msg)
	},
}

func roundDuration(dur string) string {
	arr := strings.Split(dur, " ")
	if len(arr) > 2 {
		return strings.Join(arr[:4], " ")
	}
	return strings.Join(arr[:2], " ")
}

var topTrig = regexp.MustCompile(`(?i).*!+(?:top|masters?|masterminds?|echoline).*`)
var top = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && topTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
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
			bot.Crit("!top", "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
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
		bot.Reply(m, out)
	},
}

var logmsg = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		logCTR.Lock()
		fmt.Fprintf(logCTR.File, "[%s] <%s>\t%s\n", time.Now().UTC().Format("06:01:02|15:04:05"), m.Name, m.Content)
		logCTR.Unlock()
	},
}

var logmsgBolt = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		err := setLogMSG(&Message{
			Time:    time.Now(),
			Nick:    m.Name,
			Message: m.Content,
		})
		if err != nil {
			bot.Crit("setLogMSG", "error", err)
		}
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
