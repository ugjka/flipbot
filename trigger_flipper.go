package main

import (
	"fmt"
	"regexp"

	kitty "github.com/ugjka/kittybot"
)

// This trigger flips the table

var fliptextTrig = regexp.MustCompile(`(?i)^\s*!+flip+\w*\s+(\S.*)$`)
var fliptext = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && fliptextTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, fmt.Sprintf("(╯‵Д′)╯彡%s", upside(fliptextTrig.FindStringSubmatch(m.Content)[1])))
	},
}

var unfliptextTrig = regexp.MustCompile(`(?i)^\s*!+unflip+\w*\s+(\S.*)$`)
var unfliptext = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && unfliptextTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, fmt.Sprintf("%s <(•_•<)", upside(unfliptextTrig.FindStringSubmatch(m.Content)[1])))
	},
}

func upside(input string) (output string) {
	for _, v := range input {
		if char, ok := alphabet[v]; ok {
			output = string(char) + output
		} else {
			output = string(v) + output
		}
	}
	return
}

var alphabet = map[rune]rune{
	'a':  'ɐ',
	'b':  'q',
	'c':  'ɔ',
	'd':  'p',
	'e':  'ǝ',
	'f':  'ɟ',
	'g':  'ƃ',
	'h':  'ɥ',
	'i':  'ᴉ',
	'j':  'ɾ',
	'k':  'ʞ',
	'l':  '˥',
	'm':  'ɯ',
	'n':  'u',
	'o':  'o',
	'p':  'd',
	'q':  'b',
	'r':  'ɹ',
	's':  's',
	't':  'ʇ',
	'u':  'n',
	'v':  'ʌ',
	'w':  'ʍ',
	'x':  'x',
	'y':  'ʎ',
	'z':  'z',
	'A':  '∀',
	'B':  'q',
	'C':  'Ɔ',
	'D':  'p',
	'E':  'Ǝ',
	'F':  'Ⅎ',
	'G':  'פ',
	'H':  'H',
	'I':  'I',
	'J':  'ſ',
	'K':  'ʞ',
	'L':  '˥',
	'M':  'W',
	'N':  'N',
	'O':  'O',
	'P':  'Ԁ',
	'Q':  'Q',
	'R':  'ɹ',
	'S':  'S',
	'T':  '┴',
	'U':  '∩',
	'V':  'Λ',
	'W':  'M',
	'X':  'X',
	'Y':  '⅄',
	'Z':  'Z',
	'\'': ',',
	'"':  ',',
	',':  '\'',
	'.':  '˙',
	'!':  '¡',
	'?':  '¿',
	'_':  '‾',
	'1':  'Ɩ',
	'2':  'ᄅ',
	'3':  'Ɛ',
	'4':  'ㄣ',
	'5':  'ϛ',
	'6':  '9',
	'7':  'ㄥ',
	'8':  '8',
	'9':  '6',
	'0':  '0',
	'ɐ':  'a',
	'ɔ':  'c',
	'ǝ':  'e',
	'ɟ':  'f',
	'ƃ':  'g',
	'ɥ':  'h',
	'ᴉ':  'i',
	'ɾ':  'j',
	'ʞ':  'k',
	'˥':  'l',
	'ɯ':  'm',
	'ɹ':  'r',
	'ʇ':  't',
	'ʌ':  'v',
	'ʍ':  'w',
	'ʎ':  'y',
	'∀':  'A',
	'Ɔ':  'C',
	'Ǝ':  'E',
	'Ⅎ':  'F',
	'פ':  'G',
	'ſ':  'J',
	'Ԁ':  'P',
	'┴':  'T',
	'∩':  'U',
	'Λ':  'V',
	'⅄':  'Y',
	'˙':  '.',
	'¡':  '!',
	'¿':  '?',
	'‾':  '_',
	'Ɩ':  '1',
	'ᄅ':  '2',
	'Ɛ':  '3',
	'ㄣ':  '4',
	'ϛ':  '5',
	'ㄥ':  '7',
}
