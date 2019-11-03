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

	"github.com/hako/durafmt"
	hbot "github.com/ugjka/hellabot"
	"github.com/ugjka/reverse"
	log "gopkg.in/inconshreveable/log15.v2"
)

var logCTR = struct {
	*os.File
	sync.Mutex
}{}

var seenCTR = struct {
	db map[string]*Seen
	*os.File
	*time.Ticker
	sync.RWMutex
}{
	db:     make(map[string]*Seen),
	Ticker: time.NewTicker(time.Minute * 5),
}

//Seen struct
type Seen struct {
	Seen    time.Time
	LastMSG string
}

var namesCall sync.Once

var names = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		namesCall.Do(func() {
			log.Info("firstrun", "action", "getting names")
			bot.Send("NAMES " + ircChannel)
		})
		return false
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		return false
	},
}

const seenTrig = "!seen "

var seen = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, seenTrig)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		v := strings.TrimPrefix(m.Content, seenTrig)
		if v == "" {
			return false
		}
		if strings.ToLower(v) == strings.ToLower(m.Name) {
			irc.Reply(m, fmt.Sprintf("%s: I'm seeing you!", m.Name))
			return false
		}
		seenCTR.RLock()
		if k, ok := seenCTR.db[strings.ToLower(v)]; ok {
			dur := durafmt.Parse(time.Now().UTC().Sub(k.Seen))
			if k.LastMSG != "" {
				irc.Reply(m, fmt.Sprintf("%s: I saw %s %s ago. Their last message was: %s", m.Name, v, roundDuration(dur.String()), k.LastMSG))
			} else {
				irc.Reply(m, fmt.Sprintf("%s: I saw %s %s ago", m.Name, v, roundDuration(dur.String())))
			}
		} else {
			irc.Reply(m, fmt.Sprintf("%s: I haven't seen that nick before", m.Name))
		}
		seenCTR.RUnlock()

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
		name := m.Name
		if m.Command == "KICK" {
			name = m.Params[1]
		}
		seenCTR.Lock()
		defer seenCTR.Unlock()
		if m.Command == "PRIVMSG" {
			seen, ok := seenCTR.db[strings.ToLower(name)]
			if !ok {
				seenCTR.db[strings.ToLower(name)] = &Seen{
					Seen:    time.Now().UTC(),
					LastMSG: m.Content,
				}
			} else {
				seen.Seen = time.Now().UTC()
				seen.LastMSG = m.Content
			}
		} else {
			seen, ok := seenCTR.db[strings.ToLower(name)]
			if !ok {
				seenCTR.db[strings.ToLower(name)] = &Seen{
					Seen: time.Now().UTC(),
				}
			} else {
				seen.Seen = time.Now().UTC()
			}
		}
		tmp, err := json.Marshal(seenCTR.db)
		if err == nil {
			select {
			case <-seenCTR.C:
				seenCTR.Truncate(0)
				seenCTR.WriteAt(tmp, 0)
			default:
				return false
			}
		}
		return false
	},
}

const topTrig = "!top"

var top = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.Content == topTrig
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
