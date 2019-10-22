package subwatch

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "gopkg.in/inconshreveable/log15.v2"

	"github.com/martinlindhe/base36"

	"github.com/mmcdole/gofeed"
)

// Client let's you fiddle with http.Client
var Client = &http.Client{}

type bot struct {
	endpoints      []string
	lastID         uint64
	send           chan string
	feed           *gofeed.Parser
	useragent      string
	interval       time.Duration
	round          time.Duration
	printSubreddit bool
}

//Bot settings
type Bot struct {
	Endpoints      []string
	FetchInterval  time.Duration
	Round          time.Duration
	UserAgent      string
	PrintSubreddit bool
}

//New creates a new bot object
func New(b *Bot) (reddit *bot, sender chan string) {
	msg := make(chan string, 100)
	return &bot{
			send:           msg,
			feed:           gofeed.NewParser(),
			endpoints:      b.Endpoints,
			useragent:      b.UserAgent,
			interval:       b.FetchInterval,
			round:          b.Round,
			printSubreddit: b.PrintSubreddit,
		},
		msg
}

// Get posts
func (b *bot) fetch(endpoint string) (p *gofeed.Feed, err error) {
	if !strings.Contains(endpoint, ".rss") {
		endpoint += ".rss"
	}
	req, err := http.NewRequest("GET", "https://www.reddit.com"+endpoint, nil)
	if err != nil {
		return
	}
	// Headers.
	req.Header.Set("User-Agent", b.useragent)

	resp, err := Client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("fetch response error: " + resp.Status)
	}
	return b.feed.Parse(resp.Body)
}

func (b *bot) firstRun() error {
	for _, v := range b.endpoints {
		posts, err := b.fetch(v)
		if err != nil {
			log.Info("subwatch first run", "error", err)
			return err
		}
		for _, v := range posts.Items {
			if !strings.HasPrefix(v.GUID, "t3_") {
				continue
			}
			decoded := base36.Decode(v.GUID[3:])
			if b.lastID < decoded {
				b.lastID = decoded
			}
		}
	}
	return nil
}

func (b *bot) getPosts() {
	reddit := "reddit"
	var tmpLargest uint64
	dup := make(map[uint64]bool)
	for _, v := range b.endpoints {
		posts, err := b.fetch(v)
		if err != nil {
			log.Warn("subwatch could not fetch posts", "error", err)
			return
		}
		for _, v := range posts.Items {
			if !strings.HasPrefix(v.GUID, "t3_") {
				continue
			}
			decoded := base36.Decode(v.GUID[3:])
			if _, ok := dup[decoded]; ok {
				continue
			}
			dup[decoded] = true
			if tmpLargest < decoded {
				tmpLargest = decoded
			}
			if b.lastID < decoded {
				name := ""
				if v.Author == nil {
					name = "account_deleted"
				} else {
					name = v.Author.Name
				}
				if b.printSubreddit && v.Categories != nil {
					reddit = "/r/" + v.Categories[0]
				}
				b.send <- fmt.Sprintf("[%s] [%s] %s https://redd.it/%s", reddit, name, v.Title, v.GUID[3:])
			}
		}
	}
	b.lastID = tmpLargest
}

func (b *bot) mainLoop() {
	round := time.Now().Round(b.round)
	if time.Now().After(round) {
		round = round.Add(b.round)
	}
	time.Sleep(round.Sub(time.Now()))
	ticker := time.NewTicker(b.interval)
	b.getPosts()
	for {
		select {
		case <-ticker.C:
			b.getPosts()
		}
	}
}

//Start starts the bot
func (b *bot) Start() {
	var err error
	for {
		err = b.firstRun()
		if err == nil {
			log.Info("subwatch", "success", "first run succeeded")
			break
		}
		log.Warn("subwatch", "first run failed:", err)
		time.Sleep(time.Minute * 10)
		log.Info("subwatch", "status", "retrying first run")
	}
	b.mainLoop()
}
