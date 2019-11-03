package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/ugjka/go-tz.v2/tz"
)

var clockTrig = "!time "

var clock = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, clockTrig)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		timez, err := getTime(strings.TrimPrefix(m.Content, clockTrig))
		if err != nil {
			log.Warn("could not get time", "for", m.Content[6:len(m.Content)], "error", err.Error())
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, timez))
		return false
	},
}

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
		return "I don't know that place.", nil
	}
	adress := mapj[0].DisplayName
	lat, _ := strconv.ParseFloat(mapj[0].Lat, 64)
	lon, _ := strconv.ParseFloat(mapj[0].Lon, 64)
	p := tz.Point{Lat: lat, Lon: lon}
	tzid, err := tz.GetZone(p)
	if err != nil {
		return "I don't know that place.", nil
	}
	zone, err := time.LoadLocation(tzid[0])
	if err != nil {
		return "I don't know that place.", nil
	}
	timeX := time.Now().In(zone)
	return fmt.Sprintf("Time for %s is %s", adress, timeX.Format("Mon Jan 2 15:04:05")), nil
}
