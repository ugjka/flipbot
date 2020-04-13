package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var ukcovidReg = regexp.MustCompile(`(?i)^\s*!+(?:uk|gb)(?:covid\w*)?(?:\s+(\S.*))?$`)
var ukcovid = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return ukcovidReg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		input := ukcovidReg.FindStringSubmatch(m.Content)[1]
		input = strings.TrimSpace(input)
		if input == "" {
			out, err := getUKRegion("")
			if err != nil {
				log.Error("ukcovid", "error", err)
				irc.Reply(m, fmt.Sprintf("%s: some error happened", m.Name))
				return false
			}
			irc.Reply(m, out)
			return false
		}
		if strings.Contains(input, "death") {
			out, err := getUKDeaths()
			if err != nil {
				log.Error("ukcovid", "error", err)
				irc.Reply(m, fmt.Sprintf("%s: some error happened", m.Name))
				return false
			}
			irc.Reply(m, out)
			return false
		}
		out, err := getUKRegion(input)
		if err != nil {
			log.Error("ukcovid", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: some error happened", m.Name))
			return false
		}
		irc.Reply(m, out)
		return false
	},
}

func getUKDeaths() (out string, err error) {
	const api = "https://api.covidlive.co.uk/ukdeaths.json"
	resp, err := httpClient.Get(api)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	data := make(ukDeaths, 0)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return out, err
	}
	if len(data) < 2 {
		return out, errors.New("borked data")
	}
	sorter := make([]string, 0)
	for k := range data[0].Regions {
		sorter = append(sorter, k)
	}
	sort.Strings(sorter)
	out += fmt.Sprintf("[Deaths as of %s] ", data[0].Day)
	for _, v := range sorter {
		out += fmt.Sprintf("%s: %d (+%d today), ", v, data[0].Regions[v], data[0].Regions[v]-data[1].Regions[v])
	}
	out = strings.TrimSuffix(out, ", ")
	return out, err
}

func getUKRegion(s string) (out string, err error) {
	s = strings.TrimSpace(s)
	const api = "https://api.covidlive.co.uk/ukdata.json"
	resp, err := httpClient.Get(api)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	data := make(ukConfirmed, 0)
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return out, err
	}
	if len(data) < 2 {
		return out, errors.New("borked data")
	}
	for i, v := range data {
		fmt.Sscanf(v.TotalCases, "%d", &data[i].TotalCasesInt)
		fmt.Sscanf(v.TotalDeaths, "%d", &data[i].TotalDeathsInt)
		fmt.Sscanf(v.TotalTested, "%d", &data[i].TotalTestedInt)
	}
	if s == "" {
		out += fmt.Sprintf("[UK as of %s] ", data[0].Day)
		out += fmt.Sprintf("Total cases: %d (+%d today), ", data[0].TotalCasesInt, data[0].TotalCasesInt-data[1].TotalCasesInt)
		out += fmt.Sprintf("Total tested: %d (+%d today), ", data[0].TotalTestedInt, data[0].TotalTestedInt-data[1].TotalTestedInt)
		out += fmt.Sprintf("Total deaths: %d (+%d today)", data[0].TotalDeathsInt, data[0].TotalDeathsInt-data[1].TotalDeathsInt)
		return out, nil
	}
	sorter := make([]string, 0)
	for k := range data[0].Regions {
		sorter = append(sorter, k)
	}
	sort.Strings(sorter)
	for _, v := range sorter {
		if strings.Contains(strings.ToLower(v), strings.ToLower(s)) {
			out += fmt.Sprintf("[%s as of %s] ", v, data[0].Day)
			out += fmt.Sprintf("Cases: %d (+%d today)", data[0].Regions[v], data[0].Regions[v]-data[1].Regions[v])
			return out, nil
		}
	}
	return "Region not found, see https://covidlive.co.uk/", nil
}

type ukDeaths []struct {
	Day     string
	Regions map[string]int
}

type ukConfirmed []struct {
	Day            string
	Regions        map[string]int
	TotalCases     string
	TotalCasesInt  int
	TotalTested    string
	TotalTestedInt int
	TotalDeaths    string
	TotalDeathsInt int
}
