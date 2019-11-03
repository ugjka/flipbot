package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/PuerkitoBio/goquery"
	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

const duckerTrig = "!ducker "

var ducker = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, duckerTrig)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		query := strings.TrimPrefix(m.Content, duckerTrig)
		res, err := duck(query)
		if err != nil {
			log.Warn("could not get ddg result", "for", query, "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return false
		}
		msg := fmt.Sprintf("%s: %s", m.Name, res)
		if len(msg) > textLimit {
			msg = msg[:textLimit] + "..."
		}
		irc.Reply(m, msg)
		return false
	},
}

const googleTrig = "!google "

var google = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, googleTrig)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		query := strings.TrimPrefix(m.Content, googleTrig)
		res, err := googleStuff(query)
		if err != nil {
			log.Warn("could not get google result", "for", query, "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return false
		}
		if len(res) == 0 {
			irc.Reply(m, fmt.Sprintf("%s: no results!", m.Name))
			return false
		}
		msg := fmt.Sprintf("%s: %s [%s] (%s)", m.Name, res[0].URL, res[0].Title, res[0].Abstract)
		if len(msg) > textLimit {
			msg = msg[:textLimit] + "..."
		}
		irc.Reply(m, msg)
		return false
	},
}

var googleNewsTrig = "!news "

var googlenews = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, googleNewsTrig)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		query := strings.TrimPrefix(m.Content, googleNewsTrig)
		res, err := googleNews(query)
		if err != nil {
			log.Warn("could not get google result", "for", query, "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return false
		}
		if len(res) == 0 {
			irc.Reply(m, fmt.Sprintf("%s: no results!", m.Name))
			return false
		}
		msg := fmt.Sprintf("%s: %s [%s] [%s] (%s)", m.Name,
			res[0].URL, res[0].Metadata, res[0].Title, res[0].Abstract)
		if len(msg) > textLimit {
			msg = msg[:textLimit] + "..."
		}
		irc.Reply(m, msg)
		return false
	},
}

// GooglerResult is a single result
type GooglerResult struct {
	Abstract string `json:"abstract"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

// GooglerResults is many results
type GooglerResults []GooglerResult

// GooglerNewsResult is news result
type GooglerNewsResult struct {
	Abstract string `json:"abstract"`
	Title    string `json:"title"`
	URL      string `json:"url"`
	Metadata string `json:"metadata"`
}

// GooglerNewsResults is slice of news
type GooglerNewsResults []GooglerNewsResult

func googleStuff(q string) (res GooglerResults, err error) {
	cmd := exec.Command("googler", "--lang=en", "--json", "--count=5", fmt.Sprintf("%s", q))
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return
	}
	return
}

func googleNews(q string) (res GooglerNewsResults, err error) {
	cmd := exec.Command("googler", "--lang=en", "--news", "--json", "--count=5", fmt.Sprintf("%s", q))
	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, &res); err != nil {
		return
	}
	return
}

func duck(s string) (out string, err error) {
	m := url.Values{}
	m.Add("q", s)
	req, err := http.NewRequest("GET", "https://duckduckgo.com/html/?"+m.Encode(), nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/64.0.3282.186 Safari/537.36")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	res, err := httpClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return
	}
	sel := doc.Find(".result").Not(".result--ad").First().Find(".result__a")
	url, ok := sel.Attr("href")
	if !ok {
		return "no results!", nil
	}
	return fmt.Sprintf("[%s] %s", sel.Text(), url), nil
}
