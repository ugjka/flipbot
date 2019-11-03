package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	hbot "github.com/ugjka/hellabot"
)

const transTrig = "!trans "

var trans = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, transTrig)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		res, err := translate(strings.TrimPrefix(m.Content, transTrig))
		if err != nil {
			irc.Reply(m, fmt.Sprintf("%s: %v", m.Name, err))
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
		res, err := exec.Command("trans", "-e", "google", "-brief", lang[1], lang[2]).Output()
		return strings.Replace(string(res), "\n", " ", -1), err
	}
	resByte, err := exec.Command("trans", "-e", "google", "-brief", command).Output()
	return strings.Replace(string(resByte), "\n", " ", -1), err
}
