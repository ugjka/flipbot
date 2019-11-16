package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/hako/durafmt"
	hbot "github.com/ugjka/hellabot"
	"github.com/ugjka/reverse"
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
		return m.Command == "PRIVMSG" || m.Command == "JOIN" ||
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

var topTrig = regexp.MustCompile(`(?i).*!+(?:top|masters?|masterminds?|ranks?).*`)
var top = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && topTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		logCTR.Lock()
		defer logCTR.Unlock()
		reg := regexp.MustCompile("^\\[(\\d{2}\\:\\d{2}\\:\\d{2}\\|\\d{2}\\:\\d{2}\\:\\d{2})\\] \\<(.+)\\>\\t.*$")
		stats := make(map[string]int)
		scan := reverse.NewScanner(logCTR.File)
		res := make(result, 0)
		week := time.Now().UTC().Add(time.Hour * -24 * 7)
		for scan.Scan() {
			matches := reg.FindStringSubmatch(strings.ToLower(scan.Text()))
			if len(matches) != 3 {
				continue
			}
			timestamp, err := time.Parse("06:01:02|15:04:05", matches[1])
			if err != nil {
				continue
			}
			if timestamp.After(week) {
				stats[matches[2]]++
			} else {
				break
			}
		}
		for k, v := range stats {
			k = k[:len(k)-1] + "*"
			res = append(res, stat{k, v})
		}
		sort.Sort(sort.Reverse(res))
		out := "Top posters for past 7 days: "
		for i, v := range res {
			if (i == 10) || (i == len(res)-1) {
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
