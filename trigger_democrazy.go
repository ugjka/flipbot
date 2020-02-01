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
		return m.Command == "PRIVMSG" && upvoteTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		word := upvoteTrig.FindStringSubmatch(m.Content)[1]
		word = strings.ToLower(word)
		votes, err := setvote(1, word)
		if err != nil {
			log.Crit("!upvote", "error", err)
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %.4f votes for %s. Your vote will gradually expire in 7 days",
			m.Name, votes, word))
		return false
	},
}

var downvoteTrig = regexp.MustCompile(`(?i)^\s*(?:-+|!+(?:down+|downvote+)\s+)(\S+)$`)
var downvote = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && downvoteTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		word := downvoteTrig.FindStringSubmatch(m.Content)[1]
		word = strings.ToLower(word)
		votes, err := setvote(-1, word)
		if err != nil {
			log.Crit("!downvote", "error", err)
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %.4f votes for %s. Your vote will gradually expire in 7 days",
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

var ranksTrig = regexp.MustCompile(`(?i)^\s*!+rank+\s*$`)
var ranks = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && m.To == ircChannel && ranksTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		return false
	},
}
