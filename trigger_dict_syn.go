package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var dictTrig = regexp.MustCompile(`^\s*!dict\s+(\S.*)$`)
var dict = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && dictTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		cmd := exec.Command("trans", "--no-ansi", "-d", dictTrig.FindStringSubmatch(m.Content)[1])
		errBuf := bytes.NewBuffer(nil)
		cmd.Stderr = errBuf
		out, err := cmd.Output()
		if err != nil {
			irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, errRequest))
			log.Warn("!dict", "error", errBuf)
			return false
		}
		res := ""
		defs := []string{"noun", "adjective", "verb"}
		for _, v := range defs {
			reg := regexp.MustCompile(fmt.Sprintf("((%s)\\s+(.+\\.))", v))
			matches := reg.FindStringSubmatch(string(out))
			if len(matches) == 4 {
				res = res + fmt.Sprintf("%s: %s ", strings.ToUpper(matches[2]), matches[3])
			}
		}
		if len(res) > 0 {
			irc.Reply(m, fmt.Sprintf("%s: [DEFINITIONS] %s", m.Name, limit(res)))
			return false
		}
		//Synonyms
		for _, v := range defs {
			reg := regexp.MustCompile(fmt.Sprintf(`\s+((%s)\s+-\s(.+))`, v))
			matches := reg.FindStringSubmatch(string(out))
			if len(matches) == 4 {
				res = res + fmt.Sprintf("%s: %s ", strings.ToUpper(matches[2]), matches[3])
			}
		}
		if len(res) > 0 {
			irc.Reply(m, fmt.Sprintf("%s: [SYNONYMS] %s", m.Name, limit(res)))
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: no results", m.Name))
		return false
	},
}

var synTrig = regexp.MustCompile(`^\s*!syn\s+(\S.*)$`)
var syn = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return m.Command == "PRIVMSG" && synTrig.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		cmd := exec.Command("trans", "--no-ansi", "-d", synTrig.FindStringSubmatch(m.Content)[1])
		errBuf := bytes.NewBuffer(nil)
		cmd.Stderr = errBuf
		out, err := cmd.Output()
		if err != nil {
			irc.Reply(m, fmt.Sprintf("%s: %s", m.Name, errRequest))
			log.Warn("!syn", "error", errBuf)
			return false
		}
		res := ""
		defs := []string{"noun", "adjective", "verb"}
		//Synonyms
		for _, v := range defs {
			reg := regexp.MustCompile(fmt.Sprintf(`\s+((%s)\s+-\s(.+))`, v))
			matches := reg.FindStringSubmatch(string(out))
			if len(matches) == 4 {
				res = res + fmt.Sprintf("%s: %s ", strings.ToUpper(matches[2]), matches[3])
			}
		}
		if len(res) > 0 {
			irc.Reply(m, fmt.Sprintf("%s: [SYNONYMS] %s", m.Name, limit(res)))
			return false
		}
		irc.Reply(m, fmt.Sprintf("%s: no results", m.Name))
		return false
	},
}
