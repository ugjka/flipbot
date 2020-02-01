package main

import (
	"fmt"
	"regexp"
	"strings"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var upvoteTrig = regexp.MustCompile(`(?i)^\s*(?:\++|!+(?:up+|upvote+)\s+)(\S+)$`)
var upvote = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel && upvoteTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		word := upvoteTrig.FindStringSubmatch(m.Content)[1]
		word = strings.ToLower(word)
		votes, err := setvote(1, word)
		if err != nil {
			log.Crit("!upvote", "error", err)
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %.4f votes for %s. Your upvote will gradually expire in 7 days",
			m.Name, votes, word))
		return false
	},
}

var downvoteTrig = regexp.MustCompile(`(?i)^\s*(?:-+|!+(?:down+|downvote+)\s+)(\S+)$`)
var downvote = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel && downvoteTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		word := downvoteTrig.FindStringSubmatch(m.Content)[1]
		word = strings.ToLower(word)
		votes, err := setvote(-1, word)
		if err != nil {
			log.Crit("!downvote", "error", err)
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %.4f votes for %s. Your downvote will gradually expire in 7 days",
			m.Name, votes, word))
		return false
	},
}

var rankTrig = regexp.MustCompile(`(?i)^\s*(?:\?+|!+rank+\s+)(\S+)$`)
var rank = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && rankTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		word := rankTrig.FindStringSubmatch(m.Content)[1]
		word = strings.ToLower(word)
		votes, err := getvotes(word)
		if err != nil {
			log.Crit("!rank", "error", err)
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %.4f votes for %s",
			m.Name, votes, word))
		return false
	},
}

var ranksTrig = regexp.MustCompile(`(?i)^\s*!+(?:rank+s?|leader+s?|leaderboard+s?)\s*$`)
var ranks = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && ranksTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		ranks, err := getRanks()
		if err != nil {
			log.Crit("!ranks", "error", err)
			return false
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
		out = strings.TrimSuffix(out, ",")
		irc.Reply(m, out)
		return false
	},
}
