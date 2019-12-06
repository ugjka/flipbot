package main

import (
	"fmt"
	"html"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"

	hbot "github.com/ugjka/hellabot"

	"github.com/PuerkitoBio/goquery"
	log "gopkg.in/inconshreveable/log15.v2"
	"mvdan.cc/xurls/v2"
)

var urltitle = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && xurls.Strict().MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		url := xurls.Strict().FindString(m.Content)
		res, err := getPreview(url)
		if err != nil {
			log.Warn("preview", "url", url, "error", err)
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s's link: %s", m.Name, limit(res)))
		return false
	},
}

func parseYoutubeTime(input string) string {
	units := []string{"H", "M", "S"}
	results := make([]string, 0)
	for _, v := range units {
		reg := regexp.MustCompile(fmt.Sprintf(".+[A-Z](\\d+)%s", v))
		match := reg.FindStringSubmatch(input)
		if len(match) == 2 {
			if len(match[1]) != 2 {
				match[1] = "0" + match[1]
			}
			results = append(results, match[1])
		} else {
			results = append(results, "00")
		}
	}
	return strings.Join(results, ":")
}

var youtubeIDReg = regexp.MustCompile(`(?:http[s]?\://)?(?:www\.)?youtu(?:be\.com/watch\?v=|\.be/)([0-9A-Za-z_-]{11}).*`)

func printYoutubeInfo(url string) (string, error) {
	res, err := getYoutubeTitle(youtubeIDReg.FindStringSubmatch(url)[1])
	if err != nil {
		return "", err
	}
	if len(res.Items) == 0 {
		return "", fmt.Errorf("youtube no items")
	}
	if res.Items[0].Snippet.Title == "" {
		return "", fmt.Errorf("youtube no title")
	}
	result := fmt.Sprintf("[Youtube] %s | %s | %s",
		res.Items[0].Snippet.Title,
		res.Items[0].Snippet.ChannelTitle,
		parseYoutubeTime(res.Items[0].ContentDetails.Duration))
	result = html.UnescapeString(result)
	return result, nil
}

func getPreview(url string) (title string, err error) {
	if youtubeIDReg.MatchString(url) {
		return printYoutubeInfo(url)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; rv:68.0) Gecko/20100101 Firefox/68.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	jar.Clear()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status not 200")
	}
	const htmlContent = "text/html"
	content := res.Header.Get("Content-Type")
	if !strings.Contains(content, htmlContent) {
		return "", fmt.Errorf("not %s", htmlContent)
	}
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	title = doc.Find("title").First().Contents().Text()
	if !utf8.Valid([]byte(title)) {
		return "", fmt.Errorf("not utf-8 text")
	}
	title = whitespace.ReplaceAllString(title, " ")
	title = strings.TrimSpace(title)
	title = html.UnescapeString(title)
	return title, nil
}
