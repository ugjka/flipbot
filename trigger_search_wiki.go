package main

import (
	"fmt"
	"net/url"
	"strings"

	wikimedia "github.com/pmylund/go-wikimedia"
	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

const wikiTrig = "!wiki "

var wiki = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, wikiTrig)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		answer, link, err := searchWiki(strings.TrimPrefix(m.Content, wikiTrig))
		if err != nil {
			log.Warn("could not get wiki entry", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %s [%s]", m.Name, limit(answer), link))
		return false
	},
}

func searchWiki(query string) (answer, link string, err error) {
	w, err := wikimedia.New("http://en.wikipedia.org/w/api.php")
	w.Client = httpClient
	if err != nil {
		return "", "", err
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
		return "", "", err
	}
	if len(res.Query.Search) == 0 {
		return "", "", fmt.Errorf("no results")
	}

	f = url.Values{
		"action":          {"query"},
		"prop":            {"extracts"},
		"titles":          {res.Query.Search[0].Title},
		"explaintext":     {"true"},
		"exsectionformat": {"plain"},
		"exchars":         {"740"},
		"redirects":       {"true"},
	}
	res, err = w.Query(f)
	if err != nil {
		return "", "", err
	}
	for _, v := range res.Query.Pages {
		if v.PageId == 0 {
			return "", "", fmt.Errorf("no results")
		}
		text := ""
		for _, v := range strings.Split(v.Extract, "\n") {
			text += " " + v
		}
		answer = fmt.Sprintf("%s - %s", v.Title, text)
		link = fmt.Sprintf("https://en.wikipedia.org/?curid=%d", v.PageId)
		break
	}
	return
}
