package main

import (
	"regexp"
	"strings"
	"unicode/utf8"

	kitty "flipbot/kittybot"
)

var boldReg = regexp.MustCompile(`(?i)\s*!+(?:bold)\s+(\S+.*)`)
var bold = kitty.Trigger{
	Condition: func(bot *kitty.Bot, m *kitty.Message) bool {
		return boldReg.MatchString(m.Content)
	},
	Action: func(bot *kitty.Bot, m *kitty.Message) {
		text := boldReg.FindStringSubmatch(m.Content)[1]
		text = strings.ToLower(text)
		out := ""
		maxlen := bot.ReplyMaxSize(m)
		spacer := '⚬'
		var placeholder rune
		for _, v := range text {
			if r, ok := blackbold[v]; ok {
				placeholder = r
			} else {
				placeholder = spacer
			}
			if len([]byte(out))+utf8.RuneLen(placeholder) > maxlen {
				break
			}
			out += string(placeholder)
		}
		bot.Reply(m, out)
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
