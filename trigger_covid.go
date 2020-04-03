package main

import (
	"encoding/json"
	"fmt"
	"regexp"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var covidTriggerReg = regexp.MustCompile(`(?i)\s*!+(?:covid-?(?:19)?|corona(?:virus)?)\s+(\w+.*)`)
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
		irc.Reply(m, c.String())
		return false
	},
}

var covidAllTriggerReg = regexp.MustCompile(`(?i)^\s*!+(?:covid-?(?:19)?|corona(?:virus)?)$`)
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

var coronaCountryAPI = "https://corona.lmao.ninja/countries/"
var coronaAllAPI = "https://corona.lmao.ninja/all"

type covid struct {
	Message     string
	Country     string
	Cases       int
	TodayCases  int
	Deaths      int
	TodayDeaths int
	Recovered   int
	Active      int
	Critical    int
}

func (c covid) String() string {
	if c.Message != "" {
		return c.Message
	}
	return fmt.Sprintf("[%s] cases: %d (+%d today), deaths: %d (+%d today), recovered: %d, active: %d, critical: %d",
		c.Country, c.Cases, c.TodayCases, c.Deaths, c.TodayDeaths, c.Recovered, c.Active, c.Critical)
}

type covidAll struct {
	Cases             int
	Deaths            int
	Recovered         int
	Active            int
	AffectedCountries int
}

func (c covidAll) String() string {
	return fmt.Sprintf("[Global] cases: %d, deaths: %d, recovered: %d, active: %d, affected countries: %d",
		c.Cases, c.Deaths, c.Recovered, c.Active, c.AffectedCountries)
}
