package main

import (
	"fmt"
	"os"
	"regexp"

	hbot "github.com/ugjka/hellabot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var diceTrigReg = regexp.MustCompile(`(?i)^\s*!+(?:dice+|roll+)(?:\s+(\d+))?$`)
var diceTrig = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return diceTrigReg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		rolls := 0
		arr := diceTrigReg.FindStringSubmatch(m.Content)
		if len(arr) > 1 {
			_, err := fmt.Sscanf(arr[1], "%d", &rolls)
			if err != nil {
				rolls = 1
			}
		}
		if rolls > 100 {
			rolls = 100
		}
		rand, err := os.Open("/dev/urandom")
		if err != nil {
			log.Error("dice", "rand open error", err)
			return false
		}
		defer rand.Close()
		bit := make([]byte, 1)
		out := ""
		for {
			if rolls == 0 {
				break
			}
			_, err = rand.Read(bit)
			if err != nil {
				log.Error("dice", "read error", err)
				return false
			}
			if int(bit[0]) >= 1 && int(bit[0]) <= 6 {
				out += dice[int(bit[0])] + " "
				rolls--
			}
		}
		irc.Reply(m, out)
		return false
	},
}

var dice = map[int]string{
	1: "⚀",
	2: "⚁",
	3: "⚂",
	4: "⚃",
	5: "⚄",
	6: "⚅",
}
