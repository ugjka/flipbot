package main

import (
	cookiejar "bootybot/jar"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	kitty "bootybot/kittybot"

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
	openWeatherMapAPIKey = os.Getenv(openWeatherMapAPIKeyVar)
	check(openWeatherMapAPIKey, openWeatherMapAPIKeyVar)
	op = os.Getenv(opVar)
	check(op, opVar)
	serverEmail = os.Getenv(serverEmailVar)
	check(serverEmail, serverEmailVar)
	wolframAPIKey = os.Getenv(wolframAPIKeyVar)
	check(wolframAPIKey, wolframAPIKeyVar)
	discordToken = os.Getenv(discordTokenVar)
	check(discordToken, discordTokenVar)

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

	bot := kitty.NewBot(discordToken)

	go func() {
		db, err = initDB("flipbot.db")
		if err != nil {
			panic(err)
		}
	}()

	ytErrLog.File, err = os.OpenFile("ytdl_err.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer ytErrLog.Close()

	go func() {
		<-stop
		db.Close()
		os.Exit(0)
	}()
	bot.AddTrigger(notifyop)
	bot.AddTrigger(morningTrig)
	bot.AddTrigger(pooParty)
	bot.AddTrigger(kittyParty)
	bot.AddTrigger(vixey)
	bot.AddTrigger(ukcovid)
	bot.AddTrigger(idkTrig)
	bot.AddTrigger(sexTrig)
	bot.AddTrigger(ball8)
	bot.AddTrigger(bold)
	bot.AddTrigger(diceTrig)
	bot.AddTrigger(covidTrigger)
	bot.AddTrigger(covidAllTrigger)
	bot.AddTrigger(echo)
	bot.AddTrigger(bkb)
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
	bot.AddTrigger(clock)
	bot.AddTrigger(google)
	bot.AddTrigger(wiki)
	bot.AddTrigger(fliptext)
	bot.AddTrigger(unfliptext)
	bot.AddTrigger(youtube)
	bot.AddTrigger(weatherOpen)
	bot.AddTrigger(wforecastOpen)
	bot.AddTrigger(trans)
	bot.AddTrigger(help)
	bot.AddTrigger(test)
	bot.AddTrigger(googlenews)
	bot.AddTrigger(ducker)
	bot.AddTrigger(hug)
	bot.AddTrigger(debug)
	bot.AddTrigger(calc)
	bot.AddTrigger(god)
	bot.AddTrigger(randomdog)
	bot.AddTrigger(pong)
	bot.AddTrigger(meditation)
	bot.AddTrigger(dict)
	bot.AddTrigger(syn)
	//bot.AddTrigger(youtubedl)

	logHandler := log.LvlFilterHandler(log.LvlInfo, log.StdoutHandler)
	bot.Logger.SetHandler(logHandler)

	bot.Run()
	fmt.Println("Bot shutting down.")
}
