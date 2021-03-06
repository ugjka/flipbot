package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	kitty "github.com/ugjka/kittybot"
)

var dictTrig = regexp.MustCompile(`(?i)^\s*!+dict(?:ionary)?\w*\s+(\S.*)$`)
var dict = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && dictTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		cmd := exec.Command("trans", "--no-ansi", "-d", dictTrig.FindStringSubmatch(m.Content)[1])
		errBuf := bytes.NewBuffer(nil)
		cmd.Stderr = errBuf
		out, err := cmd.Output()
		if err != nil {
			bot.Reply(m, fmt.Sprintf("%s: %s", m.Name, errRequest))
			bot.Warn("!dict", "error", errBuf)
			return
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
			msg := fmt.Sprintf("%s: [DEFINITIONS] %s", m.Name, res)
			msg = limitReply(bot, m, msg, 1)
			bot.Reply(m, msg)
			return
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
			msg := fmt.Sprintf("%s: [SYNONYMS] %s", m.Name, res)
			msg = limitReply(bot, m, msg, 1)
			bot.Reply(m, msg)
			return
		}
		bot.Reply(m, fmt.Sprintf("%s: no results", m.Name))
	},
}

var synTrig = regexp.MustCompile(`(?i)^\s*!+syn(?:onyms?)?\w*\s+(\S.*)$`)
var syn = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && synTrig.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		cmd := exec.Command("trans", "--no-ansi", "-d", synTrig.FindStringSubmatch(m.Content)[1])
		errBuf := bytes.NewBuffer(nil)
		cmd.Stderr = errBuf
		out, err := cmd.Output()
		if err != nil {
			bot.Reply(m, fmt.Sprintf("%s: %s", m.Name, errRequest))
			bot.Warn("!syn", "error", errBuf)
			return
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
			msg := fmt.Sprintf("%s: [SYNONYMS] %s", m.Name, res)
			msg = limitReply(bot, m, msg, 1)
			bot.Reply(m, msg)
			return
		}
		bot.Reply(m, fmt.Sprintf("%s: no results", m.Name))
	},
}
