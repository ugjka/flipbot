package main

import (
	"fmt"
	"os"
	"regexp"

	kitty "github.com/ugjka/kittybot"
)

var diceTrigReg = regexp.MustCompile(`(?i)^\s*!+(?:dice+|roll+)(?:\s+(\d+))?$`)
var diceTrig = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return diceTrigReg.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
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
			bot.Error("dice", "rand open error", err)
			return
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
				bot.Error("dice", "read error", err)
				return
			}
			if int(bit[0]) >= 1 && int(bit[0]) <= 6 {
				out += dice[int(bit[0])] + " "
				rolls--
			}
		}
		bot.Reply(m, out)
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
