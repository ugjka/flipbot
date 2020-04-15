package main

import (
	"regexp"
	"strings"

	hbot "github.com/ugjka/hellabot"
)

var boldReg = regexp.MustCompile(`(?i)\s*!+(?:bold)\s+(\w+.*)`)
var bold = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return boldReg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		text := boldReg.FindStringSubmatch(m.Content)[1]
		text = strings.ToLower(text)
		out := ""
		for _, v := range text {
			if r, ok := blackbold[v]; ok {
				out += string(r)
			} else {
				out += "⚬"
			}
		}
		if len(out) > 300 {
			out = out[:300]
			out = strings.ToValidUTF8(out, "")
		}
		irc.Reply(m, out)
		return false
	},
}

var blackbold = map[rune]rune{
	'a': '🅐',
	'b': '🅑',
	'c': '🅒',
	'd': '🅓',
	'e': '🅔',
	'f': '🅕',
	'g': '🅖',
	'h': '🅗',
	'i': '🅘',
	'j': '🅙',
	'k': '🅚',
	'l': '🅛',
	'm': '🅜',
	'n': '🅝',
	'o': '🅞',
	'p': '🅟',
	'q': '🅠',
	'r': '🅡',
	's': '🅢',
	't': '🅣',
	'u': '🅤',
	'v': '🅥',
	'w': '🅦',
	'x': '🅧',
	'y': '🅨',
	'z': '🅩',
	'0': '🄌',
	'1': '➊',
	'2': '➋',
	'3': '➌',
	'4': '➍',
	'5': '➎',
	'6': '➏',
	'7': '➐',
	'8': '➑',
	'9': '➒',
}
