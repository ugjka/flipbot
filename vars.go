package main

import (
	cookiejar "flipbot/jar"
	kitty "flipbot/kittybot"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/boltdb/bolt"
)

const emailVar = "FLIPBOT_EMAIL"
const youtubeAPIKeyVar = "FLIPBOT_YOUTUBE"
const subredditVar = "FLIPBOT_SUB"
const discordTokenVar = "DISCORD"
const openWeatherMapAPIKeyVar = "FLIPBOT_OW"
const opVar = "FLIPBOT_OP"
const serverEmailVar = "FLIPBOT_SERVER_MAIL"
const wolframAPIKeyVar = "FLIPBOT_WOLF"

var email string
var youtubeAPIKey string
var subreddit string
var openWeatherMapAPIKey string
var op string
var serverEmail string
var wolframAPIKey string
var discordToken string

var db = new(bolt.DB)

//Default for all requests
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

var jar *cookiejar.Jar
var meditations []string

var whitespace = regexp.MustCompile(`\s+`)

func limitReply(b *kitty.Bot, m *kitty.Message, msg string, msgCount int) string {
	limit := 512
	limit *= msgCount
	if len(msg) > limit {
		msg = msg[:limit-3] + "..."
	}
	return msg
}

var errRequest = fmt.Errorf("an error occurred while processing your request")
var errNotSeen = fmt.Errorf("nick not seen")
var errNotInCache = fmt.Errorf("not in cache")
var errNoMemo = fmt.Errorf("no memos found")
var errNoReminder = fmt.Errorf("no reminder expired")
var errNoResults = fmt.Errorf("no results")

var extJoinEnabled = false
