package main

import (
	cookiejar "flipbot/jar"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/boltdb/bolt"
	"github.com/ugjka/remindme"
)

const emailVar = "FLIPBOT_EMAIL"
const youtubeAPIKeyVar = "FLIPBOT_YOUTUBE"
const subredditVar = "FLIPBOT_SUB"
const ircServerVar = "FLIPBOT_SERVER"
const ircPortVar = "FLIPBOT_PORT"
const ircNickVar = "FLIPBOT_NICK"
const ircPasswordVar = "FLIPBOT_PASS"
const ircChannelVar = "FLIPBOT_CHAN"
const openWeatherMapAPIKeyVar = "FLIPBOT_OW"
const opVar = "FLIPBOT_OP"
const serverEmailVar = "FLIPBOT_SERVER_MAIL"
const wolframAPIKeyVar = "FLIPBOT_WOLF"

var ircPort string
var email string
var youtubeAPIKey string
var subreddit string
var ircServer string
var ircNick string
var ircPassword string
var ircChannel string
var openWeatherMapAPIKey string
var op string
var serverEmail string
var wolframAPIKey string

var remind = remindme.New(ircNick)

var db = new(bolt.DB)

//Default for all requests
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

var jar *cookiejar.Jar
var meditations []string

const textLimit = 300

var whitespace = regexp.MustCompile(`\s+`)

func limit(in string) string {
	if len(in) > textLimit {
		return in[:textLimit] + "..."
	}
	return in
}

var errRequest = fmt.Errorf("an error occurred while processing your request")
var errNotSeen = fmt.Errorf("nick not seen")
var errNotInCache = fmt.Errorf("not in cache")
var errNoMemo = fmt.Errorf("no memos found")
