package main

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	hbot "github.com/ugjka/hellabot"
)

var boldReg = regexp.MustCompile(`(?i)\s*!+(?:bold)\s+(\S+.*)`)
var bold = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return boldReg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		text := boldReg.FindStringSubmatch(m.Content)[1]
		text = strings.ToLower(text)
		out := ""
		if m.To == irc.Nick {
			m.To = m.Name
		}
		maxlen := 510 - 2 - irc.PrefixLen - len(fmt.Sprintf("PRIVMSG %s :", m.To))
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
