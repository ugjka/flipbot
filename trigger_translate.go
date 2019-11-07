package main

import (
	"fmt"
	"os/exec"
	"regexp"

	log "gopkg.in/inconshreveable/log15.v2"

	hbot "github.com/ugjka/hellabot"
)

var transTrig = regexp.MustCompile(`(?i)^\s*!+trans(?:late)?\w*\s+(\S.*)$`)
var trans = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && transTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		res, err := translate(transTrig.FindStringSubmatch(m.Content)[1])
		if err != nil {
			log.Warn("trans", "error", err)
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, limit(res)))
		return false
	},
}

func translate(command string) (res string, err error) {
	langreg := regexp.MustCompile(`^(\:[a-z]{2,3})\s+(.+)$`)
	lang := langreg.FindStringSubmatch(command)
	if len(lang) == 3 {
		out, err := exec.Command("trans", "-e", "google", "-brief", lang[1], lang[2]).Output()
		return whitespace.ReplaceAllString(string(out), " "), err
	}
	out, err := exec.Command("trans", "-e", "google", "-brief", command).Output()
	return whitespace.ReplaceAllString(string(out), " "), err
}
