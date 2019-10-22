package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	iconv "github.com/djimenez/iconv-go"
	"github.com/saintfish/chardet"
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
		suffix := regexp.MustCompile(`^https?\://.+/.+\.[a-z]{3,4}$`)
		url := xurls.Strict().FindString(m.Content)
		if suffix.MatchString(url) {
			log.Info("url is a file", "url", url)
			return false
		}
		res, err := getPreview(url)
		if err != nil {
			log.Warn("could not get url preview", "url", url, "error", err)
			return false
		}
		if res == "Too Many Requests" {
			log.Warn("could not get url preview", "url", url, "error", "Too many Requests")
			return false
		}
		res = strings.TrimSpace(res)
		if len(res) > textLimit {
			res = res[:textLimit]
		}
		irc.Reply(m, fmt.Sprintf("%s's link: %s", m.Name, res))
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

func printYoutubeInfo(url string) (string, error) {
	res, err := getYoutubeTitle(youtubeIDReg.FindStringSubmatch(url)[1])
	if err != nil {
		return "", err
	}
	if len(res.Items) == 0 {
		return "", fmt.Errorf("no items returned from youtube")
	}
	if res.Items[0].Snippet.Title == "" {
		return "", fmt.Errorf("empty title from youtube")
	}
	return fmt.Sprintf("[Youtube] %s | %s | %s",
		res.Items[0].Snippet.Title,
		res.Items[0].Snippet.ChannelTitle,
		parseYoutubeTime(res.Items[0].ContentDetails.Duration)), nil
}

func getPreview(url string) (preview string, err error) {
	if youtubeIDReg.MatchString(url) {
		preview, err = printYoutubeInfo(url)
		if err != nil {
			return
		}
	}
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:60.0) Gecko/20100101 Firefox/60.0")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	jar.Clear()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Status code not 200")
	}
	reg := regexp.MustCompile("charset=(\\w+);?")
	con := res.Header.Get("Content-Type")
	enc1 := reg.FindStringSubmatch(con)
	doc, err := goquery.NewDocumentFromResponse(res)
	if err != nil {
		return "", err
	}
	title := doc.Find("title").First().Contents().Text()
	if len(enc1) > 1 && !strings.Contains(strings.ToLower(enc1[1]), "utf") {
		tmp, err := iconv.ConvertString(title, strings.ToLower(enc1[1]), "utf-8")
		if err == nil {
			title = tmp
		}
	}
	enc, err := chardet.NewTextDetector().DetectAll([]byte(title))
	if len(enc1) < 2 && err == nil && len(enc) > 0 && !strings.Contains(strings.ToLower(enc[0].Charset), "utf") {
		tmp, err := iconv.ConvertString(title, strings.ToLower(enc[0].Charset), "utf-8")
		if err == nil {
			title = tmp
		}
	}
	preview = fmt.Sprintf("%s", title)
	if len(preview) < 4 {
		return "", fmt.Errorf("No title")
	}
	preview = strings.Replace(preview, "\n", " ", -1)
	return
}
