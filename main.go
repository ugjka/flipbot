package main

import (
	"encoding/json"
	cookiejar "flipbot/jar"
	"flipbot/subwatch"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	grmon "github.com/bcicen/grmon/agent"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

func main() {
	check := func(value, env string) {
		if value == "" {
			fmt.Fprintf(os.Stderr, "%s evnironment variable is not set\n", env)
			os.Exit(1)
		}
	}
	//Get env
	email = os.Getenv(emailVar)
	check(email, emailVar)
	youtubeAPIKey = os.Getenv(youtubeAPIKeyVar)
	check(youtubeAPIKey, youtubeAPIKeyVar)
	subreddit = os.Getenv(subredditVar)
	check(subreddit, subredditVar)
	ircServer = os.Getenv(ircChannelVar)
	check(ircServer, ircServerVar)
	ircNick = os.Getenv(ircNickVar)
	check(ircNick, ircNickVar)
	ircPassword = os.Getenv(ircPasswordVar)
	check(ircPassword, ircPasswordVar)
	ircChannel = os.Getenv(ircChannelVar)
	check(ircChannel, ircChannelVar)
	openWeatherMapAPIKey = os.Getenv(openWeatherMapAPIKeyVar)
	check(openWeatherMapAPIKey, openWeatherMapAPIKeyVar)
	op = os.Getenv(opVar)
	check(op, opVar)
	serverEmail = os.Getenv(serverEmailVar)
	check(serverEmail, serverEmailVar)
	wolframAPIKey = os.Getenv(wolframAPIKeyVar)
	check(wolframAPIKey, wolframAPIKeyVar)
	grmon.Start()

	meddata, _ := ioutil.ReadFile("meditations.txt")
	meditations = strings.Split(strings.TrimSpace(string(meddata)), "\n")
	//Cookies jar
	var err error
	jar, err = cookiejar.New(nil)
	if err != nil {
		panic(err)
	}
	httpClient.Jar = jar

	var stop = make(chan os.Signal, 3)
	signal.Notify(stop, os.Interrupt)
	signal.Notify(stop, os.Kill)
	signal.Notify(stop, syscall.SIGTERM)
	signal.Notify(stop, syscall.SIGHUP)

	hijackSession := func(bot *hbot.Bot) {
		bot.HijackSession = true
		bot.Password = ircPassword
	}

	channels := func(bot *hbot.Bot) {
		bot.Channels = []string{ircChannel}
	}
	irc, err := hbot.NewBot(ircServer, ircNick, channels, hijackSession)
	if err != nil {
		panic(err)
	}
	//Store Max online
	onlineCTR.File, err = os.OpenFile("max.json", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer onlineCTR.Close()
	err = json.NewDecoder(onlineCTR.File).Decode(&onlineCTR.Max)
	if err != nil {
		log.Warn("Cannot decode online count json file", "error", err)
	}
	//Store memo messages
	memoCTR.File, err = os.OpenFile("memo.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer memoCTR.Close()

	err = json.NewDecoder(memoCTR.File).Decode(&memoCTR.store)
	if err != nil {
		log.Warn("Cannot decode memo json file", "errro", err)
	}

	osmCTR.File, err = os.OpenFile("osmcache.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer osmCTR.Close()

	err = json.NewDecoder(osmCTR.File).Decode(&osmCTR.cache)
	if err != nil {
		log.Warn("Cannot decode osmCacheFile", "errro", err)
	}
	//Log messages
	logCTR.File, err = os.OpenFile("messages.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logCTR.Close()

	seenCTR.File, err = os.OpenFile("seen.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer seenCTR.Close()

	err = json.NewDecoder(seenCTR.File).Decode(&seenCTR.db)
	if err != nil {
		log.Warn("Cannot decode seen json file", "error", err)
	}

	go func() {
		<-stop
		logCTR.Close()
		seenCTR.Close()
		memoCTR.Close()
		onlineCTR.Close()
		os.Exit(0)
	}()

	irc.AddTrigger(flip)
	irc.AddTrigger(unflip)
	irc.AddTrigger(randomcat)
	irc.AddTrigger(ping)
	irc.AddTrigger(random)
	irc.AddTrigger(sleep)
	irc.AddTrigger(shrug)
	irc.AddTrigger(urban)
	irc.AddTrigger(define)
	irc.AddTrigger(logmsg)
	irc.AddTrigger(watcher)
	irc.AddTrigger(seen)
	irc.AddTrigger(top)
	irc.AddTrigger(clock)
	irc.AddTrigger(google)
	irc.AddTrigger(wiki)
	irc.AddTrigger(urltitle)
	irc.AddTrigger(fliptext)
	irc.AddTrigger(unfliptext)
	irc.AddTrigger(youtube)
	irc.AddTrigger(weatherOpen)
	irc.AddTrigger(wforecastOpen)
	irc.AddTrigger(trans)
	irc.AddTrigger(voice)
	irc.AddTrigger(memo)
	irc.AddTrigger(memowatcher)
	irc.AddTrigger(setmodes)
	irc.AddTrigger(voicenames)
	irc.AddTrigger(help)
	irc.AddTrigger(test)
	irc.AddTrigger(notifyop)
	irc.AddTrigger(googlenews)
	irc.AddTrigger(ducker)
	irc.AddTrigger(reminder)
	irc.AddTrigger(getreminder)
	irc.AddTrigger(onlinelist)
	irc.AddTrigger(hug)
	irc.AddTrigger(online)
	irc.AddTrigger(debug)
	irc.AddTrigger(calc)
	irc.AddTrigger(god)
	irc.AddTrigger(randomdog)
	irc.AddTrigger(pong)
	irc.AddTrigger(meditation)
	irc.AddTrigger(toss)
	irc.AddTrigger(dict)
	irc.AddTrigger(syn)
	logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	irc.Logger.SetHandler(logHandler)

	go remind.Start()

	b := &subwatch.Bot{
		Endpoints:      []string{subreddit},
		FetchInterval:  2 * time.Minute,
		Round:          2 * time.Minute,
		UserAgent:      "IRC bot for " + subreddit,
		PrintSubreddit: false,
	}
	subbot, receive := subwatch.New(b)
	go subbot.Start()
	go func() {
		for {
			irc.Msg(irc.Channels[0], <-receive)
		}
	}()
	irc.Run()
	fmt.Println("Bot shutting down.")
}
