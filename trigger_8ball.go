package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	kitty "github.com/ugjka/kittybot"
)

var ball8 = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var ball8Reg = regexp.MustCompile(`(?i)\s*!+\d*ball+(?:\s+\S*)?`)
		return ball8Reg.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		rand.Seed(time.Now().UnixNano())
		number := rand.Intn(len(ballChoices))
		b.Reply(m, fmt.Sprintf("%s: %s", m.Name, ballChoices[number]))
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
