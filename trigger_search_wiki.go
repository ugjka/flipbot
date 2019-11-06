package main

import (
	"fmt"
	"net/url"
	"regexp"

	wikimedia "github.com/pmylund/go-wikimedia"
	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var wikiTrig = regexp.MustCompile(`(?i)^\s*!wiki\s+(\S.*)$`)
var wiki = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && wikiTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		answer, link, err := searchWiki(wikiTrig.FindStringSubmatch(m.Content)[1])
		if err != nil {
			log.Warn("wiki", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
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
		err = fmt.Errorf("no results")
		return
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
		return
	}
	for _, v := range res.Query.Pages {
		if v.PageId == 0 {
			err = fmt.Errorf("no results")
			return
		}
		answer = whitespace.ReplaceAllString(v.Extract, " ")
		link = fmt.Sprintf("https://en.wikipedia.org/?curid=%d", v.PageId)
		break
	}
	return
}
