package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var covidTriggerReg = regexp.MustCompile(`(?i)\s*!+(?:covid-?(?:19)?|corona(?:virus|chan)?)\s+(\w+.*)`)
var covidTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return covidTriggerReg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		country := covidTriggerReg.FindStringSubmatch(m.Content)[1]
		resp, err := httpClient.Get(coronaCountryAPI + country)
		if err != nil {
			log.Error("covid", "get error", err)
			return false
		}
		defer resp.Body.Close()
		c := covid{}
		err = json.NewDecoder(resp.Body).Decode(&c)
		if err != nil {
			log.Error("covid", "decode error", err)
			return false
		}
		states := make(states, 0)
		state := state{}
		resp, err = httpClient.Get(coronaStatesAPI)
		if err == nil {
			defer resp.Body.Close()
			err = json.NewDecoder(resp.Body).Decode(&states)
			if err == nil {
				state = states.Search(country)
			}
		}
		switch {
		case state.State != "" && c.IsEmpty():
			irc.Reply(m, state.String())
		case state.State != "" && !c.IsEmpty():
			irc.Reply(m, c.String())
			irc.Reply(m, state.String())
		case state.State == "":
			irc.Reply(m, c.String())
		}
		return false
	},
}

var covidAllTriggerReg = regexp.MustCompile(`(?i)^\s*!+(?:covid-?(?:19)?|corona(?:virus|chan)?)$`)
var covidAllTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return covidAllTriggerReg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		resp, err := httpClient.Get(coronaAllAPI)
		if err != nil {
			log.Error("covid all", "get error", err)
			return false
		}
		defer resp.Body.Close()
		c := covidAll{}
		err = json.NewDecoder(resp.Body).Decode(&c)
		if err != nil {
			log.Error("covid all", "decode error", err)
			return false
		}
		irc.Reply(m, c.String())
		return false
	},
}

var coronaCountryAPI = "https://corona.lmao.ninja/v2/countries/"
var coronaAllAPI = "https://corona.lmao.ninja/v2/all"
var coronaStatesAPI = "https://corona.lmao.ninja/v2/states/"

type covid struct {
	Updated     int64
	Message     string
	Country     string
	Cases       int
	TodayCases  int
	Deaths      int
	TodayDeaths int
	Recovered   int
	Active      int
	Critical    int
	Tests       int
}

func (c covid) String() string {
	if c.Message != "" {
		return c.Message
	}
	return fmt.Sprintf("[%s] [%s] cases: %d (+%d today), deaths: %d (+%d today), recovered: %d, active: %d, critical: %d, tests: %d",
		time.Unix(c.Updated/1000, 0).UTC().Format("02/01/06 15:04 MST"), c.Country, c.Cases, c.TodayCases, c.Deaths, c.TodayDeaths, c.Recovered, c.Active, c.Critical, c.Tests)
}

func (c covid) IsEmpty() bool {
	return c.Message != ""
}

type covidAll struct {
	Updated           int64
	Cases             int
	TodayCases        int
	Deaths            int
	TodayDeaths       int
	Recovered         int
	Active            int
	Critical          int
	Tests             int
	AffectedCountries int
}

func (c covidAll) String() string {
	return fmt.Sprintf("[%s] [Global] cases: %d (+%d today), deaths: %d (+%d today), recovered: %d, active: %d, critical: %d, tests: %d, affected countries: %d",
		time.Unix(c.Updated/1000, 0).UTC().Format("02/01/06 15:04 MST"), c.Cases, c.TodayCases, c.Deaths, c.TodayDeaths, c.Recovered, c.Active, c.Critical, c.Tests, c.AffectedCountries)
}

type state struct {
	State       string
	Cases       int
	TodayCases  int
	Deaths      int
	TodayDeaths int
	Active      int
	Tests       int
}

type states []state

func (s states) Search(q string) state {
	q = strings.ToLower(q)
	for _, v := range s {
		if strings.ToLower(v.State) == q {
			return v
		}
	}
	return state{}
}

func (s *state) String() string {
	return fmt.Sprintf("[USA, %s] cases: %d (+%d today), deaths %d (+%d today), active: %d, tests: %d",
		s.State, s.Cases, s.TodayCases, s.Deaths, s.TodayDeaths, s.Active, s.Tests)
}
