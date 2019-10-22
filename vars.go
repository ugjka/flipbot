package main

import (
	cookiejar "flipbot/jar"
	"net/http"
	"regexp"
	"time"

	"github.com/ugjka/remindme"
)

const emailVar = "FLIPBOT_EMAIL"
const youtubeAPIKeyVar = "FLIPBOT_YOUTUBE"
const subredditVar = "FLIPBOT_SUB"
const ircServerVar = "FLIPBOT_SERVER"
const ircNickVar = "FLIPBOT_NICK"
const ircPasswordVar = "FLIPBOT_PASS"
const ircChannelVar = "FLIPBOT_CHAN"
const openWeatherMapAPIKeyVar = "FLIPBOT_OW"
const opVar = "FLIPBOT_OP"
const serverEmailVar = "FLIPBOT_SERVER_MAIL"
const wolframAPIKeyVar = "FLIPBOT_WOLF"

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

//Default for all requests
var httpClient = &http.Client{
	Timeout: 15 * time.Second,
}

var jar *cookiejar.Jar
var meditations []string

const textLimit = 300

var whitespace = regexp.MustCompile(`\s+`)
