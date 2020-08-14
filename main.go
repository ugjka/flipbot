package main

import (
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

	kitty "github.com/ugjka/kittybot"
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
	ircServer = os.Getenv(ircServerVar)
	check(ircServer, ircServerVar)
	ircPort = os.Getenv(ircPortVar)
	check(ircPort, ircPasswordVar)
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

	var err error
	meddata, err := ioutil.ReadFile("meditations.txt")
	if err != nil {
		panic(err)
	}
	meditations = strings.Split(strings.TrimSpace(string(meddata)), "\n")
	//Cookies jar
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

	hijackSession := func(b *kitty.Bot) {
		b.HijackSession = true
		b.SASL = true
		b.Password = ircPassword
	}

	channels := func(b *kitty.Bot) {
		b.Channels = []string{ircChannel}
	}
	b, err := kitty.NewBot(fmt.Sprintf("%s:%s", ircServer, ircPort), ircNick, channels, hijackSession)
	if err != nil {
		panic(err)
	}
	go func() {
		db, err = initDB("flipbot.db")
		if err != nil {
			panic(err)
		}
	}()
	//Log messages
	logCTR.File, err = os.OpenFile("messages.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logCTR.Close()

	go func() {
		<-stop
		db.Close()
		logCTR.Close()
		os.Exit(0)
	}()
	// Priority Triggers
	b.AddTrigger(logmsg)
	b.AddTrigger(logmsgBolt)
	b.AddTrigger(extJoin)
	b.AddTrigger(logJoin)
	b.AddTrigger(logAccount)
	b.AddTrigger(watcher)
	b.AddTrigger(urltitle)
	b.AddTrigger(voicenames)
	b.AddTrigger(voice)
	b.AddTrigger(names)
	b.AddTrigger(setmodes)
	b.AddTrigger(memowatcher)
	b.AddTrigger(getreminder)
	b.AddTrigger(kickmeTrigger)
	b.AddTrigger(notifyop)
	b.AddTrigger(isRecent)
	b.AddTrigger(isDead)
	// New triggers below this
	b.AddTrigger(morningTrig)
	b.AddTrigger(pooParty)
	b.AddTrigger(kittyParty)
	b.AddTrigger(vixey)
	// irc.AddTrigger(nickickerTrig)
	// irc.AddTrigger(nickickerCleanupTrig)
	b.AddTrigger(ukcovid)
	b.AddTrigger(idkTrig)
	b.AddTrigger(sexTrig)
	b.AddTrigger(ball8)
	b.AddTrigger(bold)
	b.AddTrigger(diceTrig)
	b.AddTrigger(covidTrigger)
	b.AddTrigger(covidAllTrigger)
	b.AddTrigger(upvote)
	b.AddTrigger(downvote)
	b.AddTrigger(rank)
	b.AddTrigger(ranks)
	b.AddTrigger(echo)
	b.AddTrigger(bkb)
	//irc.AddTrigger(tail)
	//irc.AddTrigger(indexUsers)
	//irc.AddTrigger(searchLog)
	//irc.AddTrigger(indexLog)
	b.AddTrigger(nature)
	b.AddTrigger(mydol)
	b.AddTrigger(flip)
	b.AddTrigger(unflip)
	b.AddTrigger(randomcat)
	b.AddTrigger(ping)
	b.AddTrigger(random)
	b.AddTrigger(sleep)
	b.AddTrigger(shrug)
	b.AddTrigger(urban)
	b.AddTrigger(define)
	b.AddTrigger(seen)
	b.AddTrigger(top)
	b.AddTrigger(clock)
	b.AddTrigger(google)
	b.AddTrigger(wiki)
	b.AddTrigger(fliptext)
	b.AddTrigger(unfliptext)
	b.AddTrigger(youtube)
	b.AddTrigger(weatherOpen)
	b.AddTrigger(wforecastOpen)
	b.AddTrigger(trans)
	b.AddTrigger(memo)
	b.AddTrigger(help)
	b.AddTrigger(test)
	b.AddTrigger(googlenews)
	b.AddTrigger(ducker)
	b.AddTrigger(reminder)
	b.AddTrigger(hug)
	b.AddTrigger(debug)
	b.AddTrigger(calc)
	b.AddTrigger(god)
	b.AddTrigger(randomdog)
	b.AddTrigger(pong)
	b.AddTrigger(meditation)
	b.AddTrigger(toss)
	b.AddTrigger(dict)
	b.AddTrigger(syn)
	// Slow triggers
	b.AddTrigger(vpnTrigger)
	b.AddTrigger(denyBETrigger)

	logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	b.Logger.SetHandler(logHandler)

	sub := &subwatch.Bot{
		Endpoints:      []string{subreddit},
		FetchInterval:  2 * time.Minute,
		Round:          2 * time.Minute,
		UserAgent:      "IRC bot for " + subreddit,
		PrintSubreddit: false,
	}
	subbot, receive := subwatch.New(sub)
	go subbot.Start()
	go func() {
		for {
			b.Msg(b.Channels[0], <-receive)
		}
	}()
	b.Run()
	fmt.Println("Bot shutting down.")
}
