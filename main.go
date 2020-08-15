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

	hijackSession := func(bot *kitty.Bot) {
		bot.HijackSession = true
		bot.SASL = true
		bot.Password = ircPassword
	}

	channels := func(bot *kitty.Bot) {
		bot.Channels = []string{ircChannel}
	}
	bot, err := kitty.NewBot(fmt.Sprintf("%s:%s", ircServer, ircPort), ircNick, channels, hijackSession)
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
	bot.AddTrigger(logmsg)
	bot.AddTrigger(logmsgBolt)
	bot.AddTrigger(extJoin)
	bot.AddTrigger(logJoin)
	bot.AddTrigger(logAccount)
	bot.AddTrigger(watcher)
	bot.AddTrigger(urltitle)
	bot.AddTrigger(voicenames)
	bot.AddTrigger(voice)
	bot.AddTrigger(names)
	bot.AddTrigger(setmodes)
	bot.AddTrigger(memowatcher)
	bot.AddTrigger(getreminder)
	bot.AddTrigger(kickmeTrigger)
	bot.AddTrigger(notifyop)
	bot.AddTrigger(isRecent)
	bot.AddTrigger(isDead)
	// New triggers below this
	bot.AddTrigger(morningTrig)
	bot.AddTrigger(pooParty)
	bot.AddTrigger(kittyParty)
	bot.AddTrigger(vixey)
	// bot.AddTrigger(nickickerTrig)
	// bot.AddTrigger(nickickerCleanupTrig)
	bot.AddTrigger(ukcovid)
	bot.AddTrigger(idkTrig)
	bot.AddTrigger(sexTrig)
	bot.AddTrigger(ball8)
	bot.AddTrigger(bold)
	bot.AddTrigger(diceTrig)
	bot.AddTrigger(covidTrigger)
	bot.AddTrigger(covidAllTrigger)
	bot.AddTrigger(upvote)
	bot.AddTrigger(downvote)
	bot.AddTrigger(rank)
	bot.AddTrigger(ranks)
	bot.AddTrigger(echo)
	bot.AddTrigger(bkb)
	//bot.AddTrigger(tail)
	//bot.AddTrigger(indexUsers)
	//bot.AddTrigger(searchLog)
	//bot.AddTrigger(indexLog)
	bot.AddTrigger(nature)
	bot.AddTrigger(mydol)
	bot.AddTrigger(flip)
	bot.AddTrigger(unflip)
	bot.AddTrigger(randomcat)
	bot.AddTrigger(ping)
	bot.AddTrigger(random)
	bot.AddTrigger(sleep)
	bot.AddTrigger(shrug)
	bot.AddTrigger(urban)
	bot.AddTrigger(define)
	bot.AddTrigger(seen)
	bot.AddTrigger(top)
	bot.AddTrigger(clock)
	bot.AddTrigger(google)
	bot.AddTrigger(wiki)
	bot.AddTrigger(fliptext)
	bot.AddTrigger(unfliptext)
	bot.AddTrigger(youtube)
	bot.AddTrigger(weatherOpen)
	bot.AddTrigger(wforecastOpen)
	bot.AddTrigger(trans)
	bot.AddTrigger(memo)
	bot.AddTrigger(help)
	bot.AddTrigger(test)
	bot.AddTrigger(googlenews)
	bot.AddTrigger(ducker)
	bot.AddTrigger(reminder)
	bot.AddTrigger(hug)
	bot.AddTrigger(debug)
	bot.AddTrigger(calc)
	bot.AddTrigger(god)
	bot.AddTrigger(randomdog)
	bot.AddTrigger(pong)
	bot.AddTrigger(meditation)
	bot.AddTrigger(toss)
	bot.AddTrigger(dict)
	bot.AddTrigger(syn)
	// Slow triggers
	bot.AddTrigger(vpnTrigger)
	bot.AddTrigger(denyBETrigger)

	logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	bot.Logger.SetHandler(logHandler)

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
			bot.Msg(bot.Channels[0], <-receive)
		}
	}()
	bot.Run()
	fmt.Println("Bot shutting down.")
}
