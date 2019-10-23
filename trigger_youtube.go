package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var youtubeIDReg = regexp.MustCompile(`(?:https?\://)?(?:www\.)?youtu(?:be\.com/watch\?v=|\.be/)([0-9A-Za-z_-]{11}).*`)

var youtube = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!youtube ")
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		res, err := searchYt(m.Content[9:])
		if err != nil {
			log.Warn("could not get youtube", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
			return false
		}
		if len(res.Items) == 0 {
			irc.Reply(m, fmt.Sprintf("%s: no results!", m.Name))
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %s https://youtu.be/%s ", m.Name, res.Items[0].Snippet.Title, res.Items[0].ID.VideoID))
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
			Title string `json:"title"`
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
