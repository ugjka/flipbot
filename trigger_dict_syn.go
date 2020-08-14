package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	kitty "github.com/ugjka/kittybot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var dictTrig = regexp.MustCompile(`(?i)^\s*!+dict(?:ionary)?\w*\s+(\S.*)$`)
var dict = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && dictTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		cmd := exec.Command("trans", "--no-ansi", "-d", dictTrig.FindStringSubmatch(m.Content)[1])
		errBuf := bytes.NewBuffer(nil)
		cmd.Stderr = errBuf
		out, err := cmd.Output()
		if err != nil {
			b.Reply(m, fmt.Sprintf("%s: %s", m.Name, errRequest))
			log.Warn("!dict", "error", errBuf)
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
			b.Reply(m, fmt.Sprintf("%s: [DEFINITIONS] %s", m.Name, limit(res, 1024)))
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
			b.Reply(m, fmt.Sprintf("%s: [SYNONYMS] %s", m.Name, limit(res, 1024)))
			return
		}
		b.Reply(m, fmt.Sprintf("%s: no results", m.Name))
	},
}

var synTrig = regexp.MustCompile(`(?i)^\s*!+syn(?:onyms?)?\w*\s+(\S.*)$`)
var syn = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && synTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		cmd := exec.Command("trans", "--no-ansi", "-d", synTrig.FindStringSubmatch(m.Content)[1])
		errBuf := bytes.NewBuffer(nil)
		cmd.Stderr = errBuf
		out, err := cmd.Output()
		if err != nil {
			b.Reply(m, fmt.Sprintf("%s: %s", m.Name, errRequest))
			log.Warn("!syn", "error", errBuf)
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
			b.Reply(m, fmt.Sprintf("%s: [SYNONYMS] %s", m.Name, limit(res, 1024)))
			return
		}
		b.Reply(m, fmt.Sprintf("%s: no results", m.Name))
	},
}
