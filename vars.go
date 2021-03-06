package main

import (
	cookiejar "flipbot/jar"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	kitty "github.com/ugjka/kittybot"
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
const mp3ServerVar = "FLIPBOT_MP3_SERVER"
const mp3DirVar = "FLIPBOT_MP3_DIR"

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
var mp3Server string
var mp3Dir string

var db = new(bolt.DB)

//Default for all requests
var httpClient = &http.Client{
	Timeout: 30 * time.Second,
}

var jar *cookiejar.Jar
var meditations []string

var whitespace = regexp.MustCompile(`\s+`)

func limitReply(b *kitty.Bot, m *kitty.Message, msg string, msgCount int) string {
	limit := b.ReplyMaxSize(m)
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

var youtubeMaxDLDur = time.Hour * 4
var youtubeDLMaxSize = "300m"

var ytErrLog = struct {
	*os.File
	*sync.Mutex
}{
	File:  new(os.File),
	Mutex: &sync.Mutex{},
}
