package main

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	kitty "github.com/ugjka/kittybot"
)

var idkTrigReg = regexp.MustCompile(`(?i)^\s*!+idk+\s+(\S.*)$`)
var idkTrig = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return idkTrigReg.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		query := idkTrigReg.FindStringSubmatch(m.Content)[1]
		link, err := idk(query)
		if err != nil {
			bot.Error("idk", "error", err)
			bot.Reply(m, fmt.Sprintf("%s: idk neither...", m.Name))
			return
		}
		bot.Reply(m, fmt.Sprintf("%s: %s", m.Name, link))
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
