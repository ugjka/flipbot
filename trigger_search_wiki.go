package main

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	kitty "flipbot/kittybot"
	wikimedia "github.com/pmylund/go-wikimedia"
)

var wikiTrig = regexp.MustCompile(`(?i)^\s*!+wiki(?:pedia)?\w*\s+(\S.*)$`)
var wiki = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && wikiTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		answer, link, err := searchWiki(wikiTrig.FindStringSubmatch(m.Content)[1])
		if err != nil {
			bot.Warn("wiki", "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errNoResults))
			return
		}
		msg := fmt.Sprintf("%s: %s", m.Name, answer)
		msg = limitReply(bot, m, msg, 5)
		bot.Reply(m, msg)
		bot.Reply(m, fmt.Sprintf("[%s]", link))
	},
}

func searchWiki(query string) (answer, link string, err error) {
	w, err := wikimedia.New("http://en.wikipedia.org/w/api.php")
	w.Client = httpClient
	if err != nil {
		return
	}
	f := url.Values{
		"action":   {"query"},
		"list":     {"search"},
		"srsearch": {query},
		"srwhat":   {"text"},
		"srprop":   {"titlesnippet"},
	}
	res, err := w.Query(f)
	if err != nil {
		return
	}
	if len(res.Query.Search) == 0 {
		link = "no results"
		return
	}

	f = url.Values{
		"action":          {"query"},
		"prop":            {"extracts"},
		"titles":          {res.Query.Search[0].Title},
		"explaintext":     {"true"},
		"exsectionformat": {"plain"},
		"exchars":         {"4096"},
		"redirects":       {"true"},
	}
	res, err = w.Query(f)
	if err != nil {
		return
	}
	for _, v := range res.Query.Pages {
		if v.PageId == 0 {
			err = fmt.Errorf("no results")
			return
		}
		answer = strings.TrimSpace(v.Extract)
		if strings.Contains(strings.Split(answer, "\n")[0], "may refer to:") {
			err = fmt.Errorf("no results")
			return
		}
		link = fmt.Sprintf("https://en.wikipedia.org/?curid=%d", v.PageId)
		break
	}
	return
}
