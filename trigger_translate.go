package main

import (
	"fmt"
	"os/exec"
	"regexp"

	kitty "bootybot/kittybot"
)

var transTrig = regexp.MustCompile(`(?i)^\s*!+trans(?:late)?\w*\s+(\S.*)$`)
var trans = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && transTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		res, err := translate(transTrig.FindStringSubmatch(m.Content)[1])
		if err != nil {
			bot.Warn("trans", "error", err)
			bot.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		msg := fmt.Sprintf("%s: %s", m.Name, res)
		msg = limitReply(bot, m, msg, 1)
		bot.Reply(m, msg)
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
