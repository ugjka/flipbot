package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	hbot "github.com/ugjka/hellabot"
)

var ball8Reg = regexp.MustCompile(`(?i)\s*!+\d*ball+(?:.*)?`)
var ball8 = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return ball8Reg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		rand.Seed(time.Now().UnixNano())
		number := rand.Intn(len(ballChoices) - 1)
		irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, ballChoices[number]))
		return false
	},
}

var ballChoices = []string{
	"Definitely",
	"Yes",
	"Probably",
	"Maybe",
	"Probably not",
	"No",
	"Definitely not",
	"I don't know",
	"Ask again later",
	"The answer is unclear",
	"Absolutely",
	"Dubious at best",
	"I'm on a break, ask again later",
	"As I see it, yes",
	"It is certain",
	"Naturally",
	"Reply hazy, try again later",
	"DO NOT WASTE MY TIME",
	"Hmm... Could be!",
	"I'm leaning towards no",
	"Without a doubt",
	"Sources say no",
	"Sources say yes",
	"Sources say maybe",
}
