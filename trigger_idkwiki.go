package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var idkTrigReg = regexp.MustCompile(`(?i)^\s*!+idk+\s+(\S.*)$`)
var idkTrig = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return idkTrigReg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		query := idkTrigReg.FindStringSubmatch(m.Content)[1]
		link, err := idk(query)
		if err != nil {
			log.Error("idk", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: idk neither...", m.Name))
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, link))
		return false
	},
}

func idk(s string) (string, error) {
	_, link, err := searchWiki(s)
	if err != nil {
		return "", err
	}
	res, err := httpClient.Get(link)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", errors.New("not 200 page")
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}
	link, _ = doc.Find("#External_links").Parent().NextUntil(".navbox").Find("ul li a.external.text[rel='nofollow']").First().Attr("href")
	if link == "" {
		return "", errors.New("no links")
	}
	return link, nil
}
