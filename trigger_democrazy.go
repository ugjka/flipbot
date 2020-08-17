package main

import (
	"fmt"
	"regexp"
	"strings"

	kitty "github.com/ugjka/kittybot"
)

var upvoteTrig = regexp.MustCompile(`(?i)^\s*(?:\++|!+(?:up+|upvote+)\s+)([[:alnum:]]\S{0,30})(?:\s+.*)?$`)
var upvote = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel && upvoteTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		word := upvoteTrig.FindStringSubmatch(m.Content)[1]
		word = strings.ToLower(word)
		votes, err := setvote(1, word)
		if err != nil {
			bot.Crit("!upvote", "error", err)
			return
		}
		bot.Reply(m, fmt.Sprintf("%s: %.4f votes for %s. Your upvote will gradually expire in 7 days",
			m.Name, votes, word))
	},
}

var downvoteTrig = regexp.MustCompile(`(?i)^\s*(?:-+|!+(?:down+|downvote+)\s+)([[:alnum:]]\S{0,30})(?:\s+.*)?$`)
var downvote = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel && downvoteTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		word := downvoteTrig.FindStringSubmatch(m.Content)[1]
		word = strings.ToLower(word)
		votes, err := setvote(-1, word)
		if err != nil {
			bot.Crit("!downvote", "error", err)
			return
		}
		bot.Reply(m, fmt.Sprintf("%s: %.4f votes for %s. Your downvote will gradually expire in 7 days",
			m.Name, votes, word))
	},
}

var rankTrig = regexp.MustCompile(`(?i)^\s*(?:\?+|!+rank+\s+)([[:alnum:]]\S{0,30})(?:\s+.*)?$`)
var rank = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && rankTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		word := rankTrig.FindStringSubmatch(m.Content)[1]
		word = strings.ToLower(word)
		votes, err := getvotes(word)
		if err != nil {
			bot.Crit("!rank", "error", err)
			return
		}
		bot.Reply(m, fmt.Sprintf("%s: %.4f votes for %s",
			m.Name, votes, word))
	},
}

var ranks = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		var ranksTrig = regexp.MustCompile(`(?i)^\s*!+(?:rank+s?|leader+s?|leaderboard+s?)\s*$`)
		return m.Command == "PRIVMSG" && ranksTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		ranks, err := getRanks()
		if err != nil {
			bot.Crit("!ranks", "error", err)
			return
		}
		out := "Leaderboard: "
		if len(ranks) == 0 {
			out += "No votes cast..."
		}
		for i, v := range ranks {
			if i > 9 {
				break
			}
			out += fmt.Sprintf("%d: %s %.4f votes, ", i+1, v.name, v.votes)
		}
		out = strings.TrimSuffix(out, ", ") + "."
		bot.Reply(m, out)
	},
}
