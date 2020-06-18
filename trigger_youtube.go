package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/url"
	"regexp"
	"time"

	"github.com/dustin/go-humanize"
	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var youtubeTrig = regexp.MustCompile(`(?i)^\s*!+(?:youtube?|yt|ytube|tube)\w*\s+(\S.*)$`)
var youtube = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && youtubeTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		res, err := searchYt(youtubeTrig.FindStringSubmatch(m.Content)[1])
		if err != nil {
			log.Warn("youtube search", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return false
		}
		if len(res.Items) == 0 {
			irc.Reply(m, fmt.Sprintf("%s: no results!", m.Name))
			return false
		}
		publishTime, err := time.Parse(time.RFC3339, res.Items[0].Snippet.PublishTime)
		if err != nil {
			log.Error("search youtube", "error", err)
			return false
		}
		result := fmt.Sprintf("%s: [Youtube] %s | %s | %s | https://youtu.be/%s",
			m.Name,
			res.Items[0].Snippet.Title,
			res.Items[0].Snippet.ChannelTitle,
			humanize.Time(publishTime),
			res.Items[0].ID,
		)
		result = html.UnescapeString(result)
		irc.Reply(m, result)
		return false
	},
}

var ytAPI = url.Values{
	"part":              {"snippet"},
	"maxResults":        {"1"},
	"type":              {"video"},
	"q":                 {""},
	"regionCode":        {"US"},
	"safeSearch":        {"none"},
	"relevanceLanguage": {"en"},
}

var ytTitleFromID = url.Values{
	"part": {"snippet,statistics,contentDetails"},
	"id":   {""},
}

type ytSearchResponse struct {
	Items []struct {
		ID struct {
			VideoID string `json:"videoId"`
		} `json:"id"`
		Snippet struct {
			Title        string `json:"title"`
			ChannelTitle string `json:"channelTitle"`
			PublishTime  string `json:"publishTime"`
		} `json:"snippet"`
	} `json:"items"`
}

type ytVideoResponse struct {
	Items []struct {
		Snippet struct {
			Title        string
			ChannelTitle string
		}
		Statistics struct {
			ViewCount    string
			LikeCount    string
			DislikeCount string
		}
		ContentDetails struct {
			Duration string
		}
	}
}

func searchYt(q string) (output ytSearchResponse, err error) {
	ytAPI.Set("q", q)
	ytAPI.Set("key", youtubeAPIKey)
	res, err := httpClient.Get("https://www.googleapis.com/youtube/v3/search?" + ytAPI.Encode())
	if err != nil {
		return
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&output); err != nil {
		return
	}
	return
}

func getYoutubeTitle(id string) (output ytVideoResponse, err error) {
	ytTitleFromID.Set("id", id)
	ytTitleFromID.Set("key", youtubeAPIKey)
	res, err := httpClient.Get("https://www.googleapis.com/youtube/v3/videos?" + ytTitleFromID.Encode())
	if err != nil {
		return
	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&output); err != nil {
		return
	}
	return
}
