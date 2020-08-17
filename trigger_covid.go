package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	kitty "github.com/ugjka/kittybot"
	"gopkg.in/ugjka/go-tz.v2/tz"
)

var covidTriggerReg = regexp.MustCompile(`(?i)\s*!+(?:covid-?(?:19)?|corona(?:virus|chan)?)\s+(\w+.*)`)
var covidTrigger = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return covidTriggerReg.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		country := covidTriggerReg.FindStringSubmatch(m.Content)[1]
		resp, err := httpClient.Get(coronaCountryAPI + country)
		if err != nil {
			bot.Error("covid", "get error", err)
			return
		}
		defer resp.Body.Close()
		c := covid{}
		err = json.NewDecoder(resp.Body).Decode(&c)
		if err != nil {
			bot.Error("covid", "decode error", err)
			return
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
			bot.Reply(m, state.String())
		case state.State != "" && !c.IsEmpty():
			bot.Reply(m, c.String())
			bot.Reply(m, state.String())
		case state.State == "":
			bot.Reply(m, c.String())
		}
	},
}

var covidAllTriggerReg = regexp.MustCompile(`(?i)^\s*!+(?:covid-?(?:19)?|corona(?:virus|chan)?)$`)
var covidAllTrigger = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return covidAllTriggerReg.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		resp, err := httpClient.Get(coronaAllAPI)
		if err != nil {
			bot.Error("covid all", "get error", err)
			return
		}
		defer resp.Body.Close()
		c := covidAll{}
		err = json.NewDecoder(resp.Body).Decode(&c)
		if err != nil {
			bot.Error("covid all", "decode error", err)
			return
		}
		bot.Reply(m, c.String())
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
	updated := time.Unix(c.Updated/1000, 0).UTC()
	if zone, err := getZone(c.Country); err == nil {
		updated = updated.In(zone)
	}
	return fmt.Sprintf("[%s] [%s] cases: %d (+%d today), deaths: %d (+%d today), recovered: %d, active: %d, critical: %d, tests: %d",
		updated.Format("02/01/06 15:04 MST"), c.Country, c.Cases, c.TodayCases, c.Deaths, c.TodayDeaths, c.Recovered, c.Active, c.Critical, c.Tests)
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

func getZone(loc string) (*time.Location, error) {
	maps := url.Values{}
	maps.Add("q", loc)
	maps.Add("format", "json")
	maps.Add("accept-language", "en")
	maps.Add("limit", "1")
	maps.Add("email", email)
	data, err := OSMGetter(OSMGeocode + maps.Encode())
	if err != nil {
		return nil, err
	}
	var mapj OSMmapResults
	if err = json.Unmarshal(data, &mapj); err != nil {
		return nil, err
	}
	if len(mapj) == 0 {
		return nil, errors.New("no place")
	}
	lat, _ := strconv.ParseFloat(mapj[0].Lat, 64)
	lon, _ := strconv.ParseFloat(mapj[0].Lon, 64)
	p := tz.Point{Lat: lat, Lon: lon}
	tzid, err := tz.GetZone(p)
	if err != nil {
		return nil, err
	}
	return time.LoadLocation(tzid[0])
}
