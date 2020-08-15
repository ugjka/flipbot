package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	kitty "github.com/ugjka/kittybot"
	log "gopkg.in/inconshreveable/log15.v2"
	"gopkg.in/ugjka/go-tz.v2/tz"
)

var nowURL = "http://api.openweathermap.org/data/2.5/weather?units=metric&lat=%s&lon=%s&APPID=%s"
var forecastURL = "http://api.openweathermap.org/data/2.5/forecast?units=metric&lat=%s&lon=%s&APPID=%s"

var errNoLocation = errors.New("location not found")

var weatherOpenTrig = regexp.MustCompile(`(?i)^\s*!+w(?:eather\w*)?\s+(\S.*)$`)
var weatherOpen = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && weatherOpenTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		lon, lat, err := getLonLat(weatherOpenTrig.FindStringSubmatch(m.Content)[1])
		if err != nil {
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		res, err := getCurrentWeather(lon, lat)
		switch err {
		case errNoLocation:
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		case nil:
			break
		default:
			log.Warn("!w", "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		format := "%s %s: %s, %.0fC/%.0fF, pressure %.1f hPa, humidity %d%%, wind %.1f m/s, gust %.1f m/s"
		a := []interface{}{
			res.Sys.Country,
			res.Name,
			res.Weather[0].Description,
			res.Main.Temp,
			res.Main.Temp*1.8 + 32,
			res.Main.Pressure,
			res.Main.Humidity,
			res.Wind.Speed,
			res.Wind.Gust,
		}
		out := fmt.Sprintf(format, a...)
		if res.Wind.Gust == 0.0 {
			format := "%s %s: %s, %.0fC/%.0fF, pressure %.1f hPa, humidity %d%%, wind %.1f m/s"
			out = fmt.Sprintf(format, a[:len(a)-1]...)
		}
		bot.Reply(m, fmt.Sprintf("%s: %s", m.Name, out))
	},
}

//OpenWNow current weather
type OpenWNow struct {
	Cod     int
	Weather []struct {
		Description string
	}
	Main struct {
		Temp     float64
		Pressure float64
		Humidity int
	}
	Wind struct {
		Speed float64
		Gust  float64
	}
	Sys struct {
		Country string
	}
	Name string
}

func getLonLat(loc string) (lon, lat string, err error) {
	maps := url.Values{}
	maps.Add("q", loc)
	maps.Add("format", "json")
	maps.Add("accept-language", "en")
	maps.Add("limit", "1")
	maps.Add("email", email)
	data, err := OSMGetter(OSMGeocode + maps.Encode())
	if err != nil {
		return
	}
	var mapj OSMmapResults
	if err = json.Unmarshal(data, &mapj); err != nil {
		return
	}
	if len(mapj) == 0 {
		return lon, lat, errNoLocation
	}
	return mapj[0].Lon, mapj[0].Lat, nil
}

func getCurrentWeather(lon, lat string) (w OpenWNow, err error) {
	resp, err := http.Get(fmt.Sprintf(nowURL, lat, lon, openWeatherMapAPIKey))
	if err != nil {
		return w, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return w, err
	}
	if w.Cod != 200 {
		return w, errNoLocation
	}
	return w, err
}

//OpenForecast is
type OpenForecast struct {
	Cod  string
	Zone string
	List []struct {
		Dt   int64
		Main struct {
			Temp        float64
			GroundLevel float64 `json:"grnd_level"`
			Humidity    int
		}
		Weather []struct {
			Description string
		}
		Wind struct {
			Speed float64
		}
	}
	City struct {
		Name    string
		Country string
	}
}

func getForecastWeather(loc string) (w OpenForecast, adress string, err error) {
	maps := url.Values{}
	maps.Add("q", loc)
	maps.Add("format", "json")
	maps.Add("accept-language", "en")
	maps.Add("limit", "1")
	maps.Add("email", email)
	data, err := OSMGetter(OSMGeocode + maps.Encode())
	if err != nil {
		return w, adress, err
	}
	var mapj OSMmapResults
	if err = json.Unmarshal(data, &mapj); err != nil {
		return w, adress, err
	}
	if len(mapj) == 0 {
		return w, adress, errNoLocation
	}
	resp, err := httpClient.Get(fmt.Sprintf(forecastURL, mapj[0].Lat, mapj[0].Lon, openWeatherMapAPIKey))
	if err != nil {
		return w, adress, err
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&w)
	if err != nil {
		return w, adress, err
	}
	if w.Cod != "200" {
		return w, adress, errNoLocation
	}
	lat, err := strconv.ParseFloat(mapj[0].Lat, 64)
	if err != nil {
		return w, adress, err
	}
	lon, err := strconv.ParseFloat(mapj[0].Lon, 64)
	if err != nil {
		return w, adress, err
	}
	zone, err := tz.GetZone(tz.Point{Lat: lat, Lon: lon})
	if err != nil {
		zone = []string{"UTC"}
	}
	w.Zone = zone[0]
	return w, adress, err
}

var wforecastOpenTrig = regexp.MustCompile(`(?i)^\s*!+(?:wf|forecast)\w*\s+(\S.*)$`)
var wforecastOpen = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && wforecastOpenTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		res, _, err := getForecastWeather(wforecastOpenTrig.FindStringSubmatch(m.Content)[1])
		switch err {
		case errNoLocation:
			bot.Reply(m, fmt.Sprintf("%s: location unknown.", m.Name))
			return
		case nil:
			break
		default:
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			log.Error("!wf", "error", err)
			return
		}
		format := "%s's forecast for %s %s"
		a := []interface{}{
			m.Name, res.City.Country, res.City.Name,
		}
		bot.Reply(m, fmt.Sprintf(format, a...))
		loc, err := time.LoadLocation(res.Zone)
		if err != nil {
			loc = time.UTC
		}
		for i, v := range res.List {
			if i%2 != 0 {
				continue
			}
			if i > 7 {
				break
			}
			currentDate := time.Unix(v.Dt, 0).In(loc)
			hour, min, _ := currentDate.Clock()
			day := currentDate.Weekday()
			format := "%s %02d:%02d %s, %.0fC/%.0fF, %.0f hPa, humidity %d%%, wind %.1f m/s"
			a := []interface{}{
				day,
				hour,
				min,
				v.Weather[0].Description,
				v.Main.Temp,
				v.Main.Temp*1.8 + 32,
				v.Main.GroundLevel,
				v.Main.Humidity,
				v.Wind.Speed,
			}
			bot.Reply(m, fmt.Sprintf(format, a...))

		}
	},
}
