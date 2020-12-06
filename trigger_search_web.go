package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	kitty "github.com/ugjka/kittybot"
)

var duckerTrig = regexp.MustCompile(`(?i)^\s*!+(?:d|(?:ducker|ddg|duck|duckduckgo)\w*)\s+(\S.*)$`)
var ducker = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && duckerTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		query := duckerTrig.FindStringSubmatch(m.Content)[1]
		res, err := duck(query)
		if err != nil {
			bot.Warn("no ducker", "for", query, "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		msg := fmt.Sprintf("%s: %s", m.Name, res)
		msg = limitReply(bot, m, msg, 1)
		bot.Reply(m, msg)
	},
}

var googleTrig = regexp.MustCompile(`(?i)^\s*!+(?:g|goog\w*)\s+(\S.*)$`)
var google = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && googleTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		query := googleTrig.FindStringSubmatch(m.Content)[1]
		res, err := googleStuff(query)
		if err != nil {
			bot.Warn("no googler", "for", query, "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		if len(res) == 0 {
			bot.Reply(m, fmt.Sprintf("%s: no results!", m.Name))
			return
		}
		msg := fmt.Sprintf("%s: %s [%s] (%s)", m.Name, res[0].URL, res[0].Title, res[0].Abstract)
		msg = limitReply(bot, m, msg, 1)
		bot.Reply(m, msg)
	},
}

var googleNewsTrig = regexp.MustCompile(`(?i)^\s*!+news+\w*\s+(\S.*)$`)
var googlenews = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && googleNewsTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		query := googleNewsTrig.FindStringSubmatch(m.Content)[1]
		res, err := googleNews(query)
		if err != nil {
			bot.Warn("no news", "for", query, "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		if len(res) == 0 {
			bot.Reply(m, fmt.Sprintf("%s: no results!", m.Name))
			return
		}
		msg := fmt.Sprintf("%s: %s [%s] [%s] (%s)", m.Name,
			res[0].URL, res[0].Metadata, res[0].Title, res[0].Abstract)
		msg = limitReply(bot, m, msg, 1)
		bot.Reply(m, msg)
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

	req, err := http.NewRequest("POST", "https://html.duckduckgo.com/html", strings.NewReader(m.Encode()))
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "IRC/Discord bot github.com/ugjka")
	req.Header.Set("Referer", "https://html.duckduckgo.com/")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
	snippet := sel.Parent().Parent().Find(".result__snippet")
	url, ok := sel.Attr("href")
	if !ok {
		return "no results!", nil
	}
	return fmt.Sprintf("%s [%s] %s", url, sel.Text(), snippet.Text()), nil
}
