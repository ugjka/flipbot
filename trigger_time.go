package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"time"

	kitty "github.com/ugjka/kittybot"
	"gopkg.in/ugjka/go-tz.v2/tz"
)

var clockTrig = regexp.MustCompile(`(?i)^\s*!+times?\w*\s+(\S.*)$`)
var clock = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && clockTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		query := clockTrig.FindStringSubmatch(m.Content)[1]
		timez, err := getTime(query)
		if err != nil {
			bot.Warn("no time", "for", query, "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		bot.Reply(m, fmt.Sprintf("%s: %s", m.Name, timez))
	},
}

const noplace = "Couldn't find such place."

//Func for querying newyears in specified location
func getTime(loc string) (string, error) {
	maps := url.Values{}
	maps.Add("q", loc)
	maps.Add("format", "json")
	maps.Add("accept-language", "en")
	maps.Add("limit", "1")
	maps.Add("email", email)
	data, err := OSMGetter(OSMGeocode + maps.Encode())
	if err != nil {
		return "", err
	}
	var mapj OSMmapResults
	if err = json.Unmarshal(data, &mapj); err != nil {
		return "", err
	}
	if len(mapj) == 0 {
		return noplace, nil
	}
	adress := mapj[0].DisplayName
	lat, _ := strconv.ParseFloat(mapj[0].Lat, 64)
	lon, _ := strconv.ParseFloat(mapj[0].Lon, 64)
	p := tz.Point{Lat: lat, Lon: lon}
	tzid, err := tz.GetZone(p)
	if err != nil {
		return noplace, nil
	}
	zone, err := time.LoadLocation(tzid[0])
	if err != nil {
		return noplace, nil
	}
	timeX := time.Now().In(zone)
	return fmt.Sprintf("Time for %s is %s", adress, timeX.Format("Mon Jan 2 15:04:05")), nil
}
