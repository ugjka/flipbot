package main

import (
	cookiejar "flipbot/jar"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	hijackSession := func(bot *hbot.Bot) {
		bot.SSL = true
	}

	channels := func(bot *hbot.Bot) {
		bot.Channels = []string{ircChannel}
	}
	irc, err := hbot.NewBot(fmt.Sprintf("%s:%s", ircServer, ircPort), ircNick, channels, hijackSession)
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
	irc.AddTrigger(pooParty)
	irc.AddTrigger(kittyParty)
	irc.AddTrigger(notifyop)
	//irc.AddTrigger(isRecent)
	//irc.AddTrigger(isDead)
	irc.AddTrigger(vixey)
	// irc.AddTrigger(nickickerTrig)
	// irc.AddTrigger(nickickerCleanupTrig)
	irc.AddTrigger(ukcovid)
	irc.AddTrigger(idkTrig)
	irc.AddTrigger(sexTrig)
	irc.AddTrigger(ball8)
	irc.AddTrigger(bold)
	irc.AddTrigger(diceTrig)
	irc.AddTrigger(covidTrigger)
	irc.AddTrigger(covidAllTrigger)
	//irc.AddTrigger(upvote)
	//irc.AddTrigger(downvote)
	//irc.AddTrigger(rank)
	//irc.AddTrigger(ranks)
	irc.AddTrigger(echo)
	irc.AddTrigger(bkb)
	//irc.AddTrigger(tail)
	//irc.AddTrigger(indexUsers)
	//irc.AddTrigger(searchLog)
	//irc.AddTrigger(indexLog)
	irc.AddTrigger(nature)
	irc.AddTrigger(mydol)
	irc.AddTrigger(flip)
	irc.AddTrigger(unflip)
	irc.AddTrigger(randomcat)
	irc.AddTrigger(ping)
	irc.AddTrigger(random)
	irc.AddTrigger(sleep)
	irc.AddTrigger(shrug)
	irc.AddTrigger(urban)
	irc.AddTrigger(define)
	//irc.AddTrigger(logmsg)
	//irc.AddTrigger(logmsgBolt)
	//irc.AddTrigger(watcher)
	//irc.AddTrigger(seen)
	//irc.AddTrigger(top)
	irc.AddTrigger(clock)
	irc.AddTrigger(google)
	irc.AddTrigger(wiki)
	//irc.AddTrigger(urltitle)
	irc.AddTrigger(fliptext)
	irc.AddTrigger(unfliptext)
	irc.AddTrigger(youtube)
	irc.AddTrigger(weatherOpen)
	irc.AddTrigger(wforecastOpen)
	irc.AddTrigger(trans)
	//irc.AddTrigger(voice)
	//irc.AddTrigger(memo)
	//irc.AddTrigger(memowatcher)
	//irc.AddTrigger(setmodes)
	//irc.AddTrigger(voicenames)
	//irc.AddTrigger(help)
	irc.AddTrigger(test)
	irc.AddTrigger(googlenews)
	irc.AddTrigger(ducker)
	//irc.AddTrigger(reminder)
	//irc.AddTrigger(getreminder)
	//irc.AddTrigger(names)
	irc.AddTrigger(hug)
	irc.AddTrigger(debug)
	irc.AddTrigger(calc)
	irc.AddTrigger(god)
	irc.AddTrigger(randomdog)
	irc.AddTrigger(pong)
	irc.AddTrigger(meditation)
	//irc.AddTrigger(toss)
	irc.AddTrigger(dict)
	irc.AddTrigger(syn)
	logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	irc.Logger.SetHandler(logHandler)

	irc.Run()
	fmt.Println("Bot shutting down.")
}
