package main

import (
	mar "flipbot/markov"
	"fmt"
	"math/rand"
	"os/exec"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/ugjka/catrand"
	kitty "github.com/ugjka/kittybot"
	log "gopkg.in/inconshreveable/log15.v2"
)

var vixey = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var vixeyTrig = regexp.MustCompile(`(?i).*!+(?:vixey.*|dash.*)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && vixeyTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		const vixeyIMG = "https://i.imgur.com/6Amydph.jpg"
		b.Reply(m, vixeyIMG)
	},
}

var morningTrig = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var morningReg = regexp.MustCompile(`(?i).*!+(?:goo+d)?morning(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && morningReg.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		morning := []rune("☕☕☕☕🥪🥐🧈🍞🥞🥓😴🥱🙃⏰💻🧇🚬🚿🚽🪥🥣")
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(morning), func(i, j int) {
			morning[i], morning[j] = morning[j], morning[i]
		})
		b.Reply(m, string(morning))
	},
}

var kittyParty = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var kittyReg = regexp.MustCompile(`(?i).*!+(?:kit+y+|kitten|cat+)party+(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && kittyReg.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		cats := []rune("😺😸😹😻😼😽🙀😿😾🐈🐈‍")
		who := m.To
		if m.To == b.Nick {
			who = m.Name
		}
		maxlen := 510 - 2 - b.Prefix().Len() - len(fmt.Sprintf("PRIVMSG %s :", who))
		msg := ""
		rand.Seed(time.Now().UnixNano())
		for i := 0; ; {
			j := rand.Intn(len(cats))
			if i+utf8.RuneLen(cats[j]) > maxlen {
				break
			}
			msg += string(cats[j])
			i += utf8.RuneLen(cats[j])
		}
		b.Reply(m, msg)
	},
}

var pooParty = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var pooReg = regexp.MustCompile(`(?i).*!+(?:po+p?|shit+y?|crap)party+(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && pooReg.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		poo := []rune("‍💩🚽🧻🍷🍺😵🤢🤮💊💉🌧️🥦🧄🍄🌿")
		who := m.To
		if m.To == b.Nick {
			who = m.Name
		}
		maxlen := 510 - 2 - b.Prefix().Len() - len(fmt.Sprintf("PRIVMSG %s :", who))
		msg := ""
		rand.Seed(time.Now().UnixNano())
		for i := 0; ; {
			j := rand.Intn(len(poo))
			if i+utf8.RuneLen(poo[j]) > maxlen {
				break
			}
			msg += string(poo[j])
			i += utf8.RuneLen(poo[j])
		}
		b.Reply(m, msg)
	},
}

var echo = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && strings.HasPrefix(m.Content, "!repeat ")
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		msg := fmt.Sprintf("%s says: '%s'", m.Name, strings.TrimPrefix(m.Content, "!repeat "))
		b.Reply(m, msg)
	},
}

var test = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var testTrig = regexp.MustCompile(`(?i).*!+(?:test|testing|check|caddy\w*|ceph\w*)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && testTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, "Congratulations, all kittens tested and ready!")
	},
}

var hugTrig = regexp.MustCompile(`(?i)^\s*!+(?:hugs?|loves?)\s+(\S.*)$`)
var hug = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		return m.Command == "PRIVMSG" && hugTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, fmt.Sprintf("%s hugs %s!", m.Name, hugTrig.FindStringSubmatch(m.Content)[1]))
	},
}

var randomdog = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var randomdogTrig = regexp.MustCompile(`(?i).*!+(?:dog+|dog+o|goodboi|pup+|pup+er|pup+ie)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && randomdogTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(len(dogs))
		b.Reply(m, dogs[n])
	},
}

var shrug = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var shrugTrig = regexp.MustCompile(`(?i).*!+(?:shrug|srug|shug|unas\w*)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && shrugTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, "¯\\_(ツ)_/¯")
	},
}

var sleep = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var sleepTrig = regexp.MustCompile(`(?i).*!+(?:sleep|nn|nite|goodnight|night|bed|nap)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && sleepTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, "【☆goodnight☆】(●ＵωU).zZZ")
	},
}

var random = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var randomTrig = regexp.MustCompile(`(?i).*!+(?:random|mad|madcotto|cotto|salad)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && randomTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, mar.Get("./for_sz_markov.txt"))
	},
}

var ping = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var pingTrig = regexp.MustCompile(`(?i).*!+ping+(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && pingTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, "PONG")
	},
}

var pong = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var pongTrig = regexp.MustCompile(`(?i).*!+pong+(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && pongTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, "PING")
	},
}

var flip = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var flipTrig = regexp.MustCompile(`^(?i).*!+(?:flip+|tableflip|fliptable)(?:\s+\S*)?$`)
		return m.Command == "PRIVMSG" && flipTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, "(╯‵Д′)╯彡┻━┻")
	},
}

var unflip = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var unflipTrig = regexp.MustCompile(`^(?i).*!+unflip+(?:\s+\S*)?$`)
		return m.Command == "PRIVMSG" && unflipTrig.MatchString(m.Content)
	},
	Action: func(irc *kitty.Bot, m *kitty.Message) {
		irc.Reply(m, "┳━┳ <(•_•<)")
	},
}

var randomcat = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var randomcatTrig = regexp.MustCompile(`^(?i)!+(?:cat+|kit+y|fluf+|kit+en+|bagpus+|pus+|pus+y)$`)
		return m.Command == "PRIVMSG" && randomcatTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, catrand.GetCat())
	},
}

var define = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var defineTrig = regexp.MustCompile(`(?i).*!+define(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && defineTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, fmt.Sprintf("%s: !define is now !urban", m.Name))
	},
}

var toss = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var tossTrig = regexp.MustCompile(`(?i).*!+(?:tos+|wank|cum+|come|shek\w*)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && tossTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		if strings.HasPrefix(m.Name, "shekib") {
			b.Reply(m, fmt.Sprintf("%s: [Hot Lebanese chick loses virginity!] https://www.youtube.com/watch?v=9y4JwyjdY4E", m.Name))
			return
		}
		text, err := tosss()
		if err != nil {
			log.Warn("!toss", "error", err)
			b.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		b.Reply(m, fmt.Sprintf("%s: %s", m.Name, text))
	},
}

var meditation = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var meditationTrig = regexp.MustCompile(`(?i).*!+(?:meditation|meditate|advaita|monism|wisdom|ugjka)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && meditationTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(len(meditations))
		b.Reply(m, fmt.Sprintf("%s: \"%s\"", m.Name, meditations[n]))
	},
}

var mydol = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var mydolTrig = regexp.MustCompile(`(?i)!+m+y+d+o+l+(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && mydolTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, "https://www.amazon.com/l/B076QJR7LF")
	},
}

var nature = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var natureTrig = regexp.MustCompile(`(?i)!nature+(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && natureTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, "https://www.flightradar24.com/")
	},
}

var god = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var godTrig = regexp.MustCompile(`(?i).*!+(?:gods?|almighty|gibberish)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && godTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		cmd := exec.Cmd{Path: "./words.sh"}
		data, err := cmd.Output()
		if err != nil {
			log.Warn("!god", "error", err)
			b.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		b.Reply(m, fmt.Sprintf("God says: %s", string(data)))
	},
}

var bkb = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var bkbTrig = regexp.MustCompile(`(?i).*!+(?:b+k+b+|e+rowid+|t+r+i+p+.*|d+r+u+g+s+)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && bkbTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		text, err := randomErowid()
		if err != nil {
			log.Warn("!bkb", "error", err)
			b.Reply(m, fmt.Sprintf("%s: %v", m.Name, errRequest))
			return
		}
		b.Reply(m, fmt.Sprintf("%s: %s", m.Name, text))
	},
}

func randomErowid() (string, error) {
	const url = "https://erowid.org/experiences/exp.php?ID=%d"
	const max = 113706
	rand.Seed(time.Now().UnixNano())
	item := rand.Int31n(max-1) + 1
	resp, err := httpClient.Get(fmt.Sprintf(url, item))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return "", err
	}
	text := doc.Find(".ts-citation").First().Text()

	if text == "" {
		return randomErowid()
	}
	text = text[14:]
	text = strings.Replace(text, "erowid.org/exp/", "https://erowid.org/exp/", 1)
	text = strings.Replace(text, " Erowid.org.", "", 1)
	fmt.Println(text)
	return text, nil
}

var help = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var helpTrig = regexp.MustCompile(`(?i)^!+(?:help|manual|com+ands|list)(?:\s+\S*)?$`)
		return m.Command == "PRIVMSG" && helpTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		b.Reply(m, "Fl1pbot's manual: https://raw.githubusercontent.com/ugjka/flipbot/bootybot/help.txt")
	},
}

var debug = kitty.Trigger{
	Condition: func(b *kitty.Bot, m *kitty.Message) bool {
		var debugTrig = regexp.MustCompile(`(?i).*!+(?:debug|bug|joke|xyk)(?:\s+\S*)?`)
		return m.Command == "PRIVMSG" && debugTrig.MatchString(m.Content)
	},
	Action: func(b *kitty.Bot, m *kitty.Message) {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(len(jokes))
		b.Reply(m, jokes[r])
	},
}

var jokes = []string{
	"unzip, strip, touch, finger, grep, mount, fsck, more, yes, fsck, fsck, fsck, umount, sleep",
	"“Knock, knock.” “Who’s there?” very long pause…. “Java.”",
	"A SQL query goes into a bar, walks up to two tables and asks, 'Can I join you?'",
	"Q: how many programmers does it take to change a light bulb? A: none, that's a hardware problem",
	"When your hammer is C++, everything begins to look like a thumb.",
	"If you put a million monkeys at a million keyboards, one of them will eventually write a Java program.",
	"Q: Whats the object-oriented way to become wealthy? A: Inheritance",
	"['hip','hip'] (hip hip array!)",
	"Programming is like sex: One mistake and you have to support it for the rest of your life.",
	"Software is like sex: It's better when it's free. (Linus Torvalds)",
	"Q: How many prolog programmers does it take to change a lightbulb? A: Yes.",
	"To understand what recursion is, you must first understand recursion.",
	"so this programmer goes out on a date with a hot chick",
	"There are 10 types of people in the world. Those who understand binary and those who have regular sex.",
	"[ $[ $RANDOM % 6 ] == 0 ] && rm -rf / || echo *Click*",
	"Unix is user friendly. It's just very particular about who its friends are.",
	"A programmer puts two glasses on his bedside table before going to sleep. A full one, in case he gets thirsty, and an empty one, in case he doesn't.",
	"A foo walks into a bar, takes a look around and says 'Hello World!' and meet up his friend Baz",
	"Q: Why don't jokes work in octal? A: Because 7 10 11.",
	"If your mom was a collection class, her insert method would be public.",
	"Female software engineers become sexually irresistible at the age of consent, and remain that way until about thirty minutes after clinical death. Longer if it's a warm day.",
	"The C language combines all the power of assembly language with all the ease-of-use of assembly language.",
	"Keyboard not found ... press F1 to continue",
	"Don't anthropomorphize computers. They hate that!",
	"Two bytes meet. The first byte asks, “Are you ill?” The second byte replies, “No, just feeling a bit off.”",
	"Specifications are for the weak and timid!",
	"You question the worthiness of my code? I should kill you where you stand!",
	"Indentation? I will show you how to indent when I indent your skull!",
	"Two threads walk into a bar. The barkeeper looks up and yells, hey, I want don't any conditions race like time last!",
	"Why doesn't C++ have a garbage collector? Because there would be nothing left!",
	"Smith & Wesson - the original 'point and click' interface.",
	"Why are Assembly programmers always soaking wet? They work below C-level.",
	"In theory, there ought to be no difference between theory and practice. In practice, there is.",
	"Nothing seems hard to the people who don't know what they're talking about.",
	"Your mommas so fat that not even Dijkstra is able to find a shortest path around her.",
	"C++ - where your friends have access to your private members.",
	"A good programmer is someone who looks both ways before crossing a one-way street. ~ Doug Linder",
	"The only 'intuitive' user interface is the nipple. After that, it's all learned.",
	"Q: Why did the programmer quit his job? A: Because he didn't get arrays.",
	"XML is like violence. If it doesn't solve your problem, you're not using enough of it",
	"Software developers like to solve problems. If there are no problems handily available, they will create their own problems.",
	"I'd like to make the world a better place, but they won't give me the source code.",
	"There's no place like 127.0.0.1",
	"I don't see women as objects. I consider each to be in a class of her own.",
	".NET is called .NET so that it wouldn't show up in a Unix directory listing",
	"What do you mean, it needs comments!? If it was hard to write, it should be hard to understand--why do you think we call it code???",
	"Hardware: The part of a computer that you can kick.",
	"Your momma's so fat, that when she sat on a binary tree she turned it into a sorted linked-list in O(1).",
	"It compiles! Let's ship it.",
	"In C we had to code our own bugs. In C++ we can inherit them.",
	"Q: How come there is not obfuscated Perl contest? A: Because everyone would win.",
	"Documentation is like sex. When it's good, it's very good. When it's bad, it's better than nothing.",
	"Q: How many programmers does it take to kill a cockroach? A: Two: one holds, the other installs Windows on it",
	"It works, don't touch!",
	"Programmers are machines that turn coffee into code.",
	"I � Unicode.",
	"Q: Why did the concurrent chicken cross the road? A: the side other To to get",
	"Walking on water and developing software from a specification are easy if both are frozen.",
	"UNIX is like eating insects. It's all right once you get used to it.",
	"Q - Why don't programmers pray? A - They don't like throwing null pointer exceptions!",
	"Some call me '^F[a-z\\'-]+$', but I have many names.",
	"Save the mallocs, free them all!",
	"Every time the God divides by zero a black hole is spawned.",
	"there's no faster code than no code!",
	"if only you and dead people can read hex, how many people can read hex?",
	"Security is not a process, it's a thread!",
	"Programming today is a race between software engineers striving to build bigger and better idiot-proof programs, and the Universe trying to produce bigger and better idiots. So far, the Universe is winning.",
	"Once you hit 10 nested VMs you're basically brad pitt with women who matter",
}

func tosss() (string, error) {
	resp, err := httpClient.Get("https://www.pornhub.com/video/random")
	if err != nil {
		return "", err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}
	if doc.Find(".premiumLocked").Text() != "" {
		return tosss()
	}
	title := doc.Find(".inlineFree").First().Contents().Text()
	return fmt.Sprintf("[%s] %s", title, resp.Request.URL), nil
}

var dogs = []string{
	"ฅ^•ﻌ•^ฅ", "੯ੁૂ‧̀͡u", "(❍ᴥ❍ʋ)", "₍ᐢ•ﻌ•ᐢ₎*･ﾟ｡", "ཥමཙමཤ", "ং” ა",
	"(^・(I)・^)", "(^・x・^)", "ｖ・。・Ｖ",
	"V●ᴥ●V", "V◕ฺω◕ฺV", "(V●ᴥ●V)",
	"∪･ω･∪", "(U・x・U) ", "┌U･ｪ･U┘",
	"ｏ（Ｕ・ω・）⊃", "U ´꓃ ` U ", "U・♀・U",
	"U｡･ｪ･｡U", "U＾ェ＾U", "Ｕ^皿^Ｕ",
	"U￣ｰ￣U", "Uo･ｪ･oU ", "ＵＴｪＴＵ",
	"(U •́ .̫ •̀ U)", "(꒪ω꒪υ)", "(υ◉ω◉υ)",
	"Uo-ｪ-oU", "(Ｕ◕ฺ㉨◕ฺ)ノ ", "(υŏᆺŏυ)",
	"ヽUﾟ●_ﾟ*Uﾉ", "U｡-ｪ-｡U。", "u(´Д`u)",
	"(Ｕ^ω^)", "(∪＾ω＾)",
	"♪o(･x･o∪ ∪o･x･)o♪",
	".+:｡ヽUﾟДﾟUﾉﾟ.+:｡",
	"(〓￣(∵エ∵)￣〓) ", "“v(〓￣(∵エ∵)￣〓)v”",
	"ヾ(〓￣(∵エ∵)￣〓) ",
	"▼o・ェ・o▼", "▽・ｗ・▽", "▽・ω・▽",
	"▿‧͈•̻‧͈▿ ", "▼･。･▼ ", "▽･ｪ･▽ﾉ”",
	"◖⚆ᴥ⚆◗ ", "⎩ ♨ᴥ♨ ⎭", "-ᄒᴥᄒ-",
	"[⑇◍ᴥ◍•⑇]", "(ノ ̿ ̿ᴥ ̿ ̿)ノ", "୧༼◕ ᴥ ◕༽୨",
	"ᕙ༼◕ ᴥ ◕༽ᕗ", "(_/¯⊘_ᴥ_⊘)_/¯ ", "⊆ↂᴥↂ⊇",
	"⎰≀.⎔ᴥ⎔≀⎰", "⊂▶ᴥ◀⊃", "°˖✧◝(ਠᴥਠ)◜✧˖°",
	"(っ⊂•⊃_ᴥ_⊂•⊃)っ", "●ᴥ●", "ヽ(°ᴥ°)ﾉ",
	"└(°ᴥ°)┘", "┏(°ᴥ°)┓", "へ║ ◉ ᴥ ◉ ║〜",
	"乁[ ◕ ᴥ ◕ ]ㄏ", "(‷\\(ᓄ ᴥ ᓇ)/‴) ", "▐ ☯ ᴥ ☯ ▐",
	"໒( ̿ ᴥ ̿ )७", "໒( ◉ ᴥ ◉ )७", "| * O ᴥ O * |",
	"٩། ಠ ᴥ ಠ །ᕗ", "୧╏ ~ ᴥ ~ ╏୨ ", "⋋〳 ￣ ᴥ ￣ 〵⋌",
	"╏ ◯ ᴥ ◯ ╏", "੧〳 ˵ ಠ ᴥ ಠ ˵ 〵ノ⌒.", "( ͡° ᴥ ͡°)",
	"…(๑╯ﻌ╰๑)=3", "୧| ⁰ ᴥ ⁰ |୨", "໒(◉ᴥ◉)७",
	"ᘳ´• ᴥ •`ᘰ", "⁞ ✿ ᵒ̌ ᴥ ᵒ̌ ✿ ⁞ ",
	"｜｡･)‐⌒ε==3 ﾍU^ｪ^U",
	"♪♪♪ Ｕ・ｪ・Ｕ人(^･x･^=) ♪♪♪",
	"o(･ω･｡)o—∈･^ミ┬┬~",
	"o(^-^ )o——–⊆^U)┬┬~",
	"o(￣_￣|||)o——–⊆◎U)┬┬ﾉ~”♪♪…",
	"o(^^ )o——–⊆^U)┬┬~…",
	"ヾ(;ﾟ皿ﾟ)ﾉ･･･ ⊆￣U)┬┬ﾉ~”　=3 =3",
	"⊂ﾟＵ┬────┬~ ", "⊂’Ｕ",
	"ε==3 ⊆＾ ￣⊇ゝ ", "∈･^ミ┬┬~",
	"⊆^U)┬┬~ ", "⊆◎U)┬┬",
	"⊆￣U)┬┬ﾉ~ ",
	"ヾ(●ω●)ノ", "(●⌇ຶ ཅ⌇ຶ●) ", "└@(･ｪ･)@┐",
	"Ψ( ◉ཅ◉ )Ψ ", "Ψ(●` ཅ ´●)Ψ ", "( ͒꒪̛ཅ꒪̛ ͒)✧",
	"(⌯꒪͒ ૢཅ ૢ꒪͒)｡*ﾟ✧",
	"Ψ(●°̥̥̥̥̥̥̥̥ ཅ °̥̥̥̥̥̥̥̥●)Ψ ", "(๑•ิཬ•ั๑) ", "༲ྕ༲",
	"(❝᷁॔Ꭳ❝᷀॓)",
}
