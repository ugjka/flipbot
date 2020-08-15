package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hako/durafmt"
	log "gopkg.in/inconshreveable/log15.v2"

	kitty "github.com/ugjka/kittybot"
)

var remindOnce sync.Once
var reminder = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.To == ircChannel
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		remindOnce.Do(func() {
			go func() {
				ticker := time.Tick(time.Second)
			Loop:
				for {
					select {
					case <-ticker:
						reminders, err := getReminder()
						switch {
						case err == errNoReminder:
							continue Loop
						case err != nil:
							log.Crit("getReminders", "error", err)
							return
						}
						for _, v := range reminders {
							text := fmt.Sprintf("%s's reminder: %s", v.Name, v.Message)
							bot.Msg(ircChannel, text)
						}
					}
				}
			}()
		})
	},
}

var getreminderTrig = regexp.MustCompile(`(?i)^\s*!+remind(?:er|me)?\w*\s+(\S.*)$`)
var getreminder = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.To == ircChannel && m.Command == "PRIVMSG" &&
			getreminderTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		target, r, err := parse(getreminderTrig.FindStringSubmatch(m.Content)[1])
		if err != nil {
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return
		}
		r.Name = m.Name
		err = setReminder(target.Format(time.RFC3339), r)
		if err != nil {
			log.Crit("setReminder", "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		delta := target.Sub(time.Now())
		dur := durafmt.Parse(delta)
		bot.Reply(m, fmt.Sprintf("%s: Your reminder will fire %s from now", m.Name, roundDuration(dur.String())))
	},
}

//ErrNoMessage no message
var errNoMessage = fmt.Errorf("%s", "No message given")

//ErrFailedParse fails parsing
var errFailedParse = fmt.Errorf("%s", "Failed to parse duration")

//ErrDurationTooBig duration past year 5000
var errDurationTooBig = fmt.Errorf("%s", "Duration too big")

func duration(d time.Duration) func(int, time.Time) time.Time {
	return func(i int, t time.Time) time.Time {
		return t.Add(time.Duration(i) * d)
	}
}

func date(y int, m int, d int) func(int, time.Time) time.Time {
	return func(i int, t time.Time) time.Time {
		return t.AddDate(y*i, m*i, d*i)
	}
}

var durations = map[string]func(int, time.Time) time.Time{
	"s":         duration(time.Second),
	"sec":       duration(time.Second),
	"second":    duration(time.Second),
	"seconds":   duration(time.Second),
	"m":         duration(time.Minute),
	"min":       duration(time.Minute),
	"mins":      duration(time.Minute),
	"minute":    duration(time.Minute),
	"minutes":   duration(time.Minute),
	"h":         duration(time.Hour),
	"hour":      duration(time.Hour),
	"hours":     duration(time.Hour),
	"d":         date(0, 0, 1),
	"day":       date(0, 0, 1),
	"days":      date(0, 0, 1),
	"tomorrow":  date(0, 0, 1),
	"w":         date(0, 0, 7),
	"week":      date(0, 0, 7),
	"weeks":     date(0, 0, 7),
	"mon":       date(0, 1, 0),
	"month":     date(0, 1, 0),
	"months":    date(0, 1, 0),
	"y":         date(1, 0, 0),
	"year":      date(1, 0, 0),
	"years":     date(1, 0, 0),
	"century":   date(100, 0, 0),
	"centuries": date(100, 0, 0),
}

//ReminderItem is a reminder
type ReminderItem struct {
	Name    string
	Message string
}

//ReminderItems is a slice
type ReminderItems []ReminderItem

var timeForm = regexp.MustCompile(`(\d+)([a-zA-Z]+)`)

//Parse reminder
func parse(m string) (target time.Time, r ReminderItem, err error) {
	m = strings.TrimSpace(m)
	m = timeForm.ReplaceAllStringFunc(m, func(in string) string {
		args := timeForm.FindStringSubmatch(in)
		if len(args) == 3 {
			return fmt.Sprintf("%s %s", args[1], args[2])
		}
		return in
	})
	arr := strings.Split(m, " ")
	for i, v := range arr {
		arr[i] = strings.Trim(v, ",")
	}

	current := now()
	i := 0
	for i < len(arr) {
		num, err := strconv.Atoi(arr[i])
		if err == nil {
			if v, ok := durations[arr[i+1]]; ok {
				current = v(num, current)
				i += 2
				continue
			}
			break
		}
		if v, ok := durations[arr[i]]; ok {
			current = v(1, current)
			i++
			continue
		}
		break
	}
	r.Name = ""
	r.Message = strings.Join(arr[i:], " ")

	if r.Message == "" {
		return now(), r, errNoMessage
	}
	if current.Year() > 5000 {
		return now(), r, errDurationTooBig
	}
	if current.Before(now()) {
		return now(), r, errFailedParse
	}

	return current, r, nil
}

//Wrapper to enable testing
var now = func() time.Time {
	return time.Now()
}
