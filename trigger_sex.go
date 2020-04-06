package main

import (
	"math/rand"
	"regexp"
	"time"

	hbot "github.com/ugjka/hellabot"
)

var sexTrigReg = regexp.MustCompile(`(?i)^\s*!+sex+$`)
var sexTrig = hbot.Trigger{
	Condition: func(bot *hbot.Bot, m *hbot.Message) bool {
		return sexTrigReg.MatchString(m.Content)
	},
	Action: func(irc *hbot.Bot, m *hbot.Message) bool {
		irc.Reply(m, altogethernow())
		return false
	},
}

var faster = []string{"\"Let the games begin!\"", "\"Sweet Jesus!\"", "\"Not that!\"", "\"At last!\"", "\"Land o' Goshen!\"", "\"Is that all?\"", "\"Cheese it, the cops!\"", "\"I never dreamed it could be\"", "\"If I do, you won't respect me!\"", "\"Now!\"", "\"Open sesame!\"", "\"EMR!\"", "\"Again!\"", "\"Faster!\"", "\"Harder!\"", "\"Help!\"", "\"Fuck me harder!\"", "\"Is it in yet?\"", "\"You aren't my father!\"", "\"Doctor, that's not *my* shou\"", "\"No, no, do the goldfish!\"", "\"Holy Batmobile, Batman!\"", "\"He's dead, he's dead!\"", "\"Take me, Robert!\"", "\"I'm a Republican!\"", "\"Put four fingers in!\"", "\"What a lover!\"", "\"Talk dirty, you pig!\"", "\"The ceiling needs painting,\"", "\"Suck harder!\"", "\"The animals will hear!\"", "\"Not in public!\""}

var said = []string{"bellowed", "yelped", "croaked", "growled", "panted", "moaned", "grunted", "laughed", "warbled", "sighed", "ejaculated", "choked", "stammered", "wheezed", "squealed", "whimpered", "salivated", "tongued", "cried", "screamed", "yelled", "said"}

var the = []string{"the"}

var fadj = []string{"saucy", "wanton", "unfortunate", "lust-crazed", "100-year-old", "bull-dyke", "bisexual", "gorgeous", "sweet", "nymphomaniacal", "large-hipped", "freckled", "forty-five year old", "white-haired", "large-boned", "saintly", "blind", "bearded", "blue-eyed", "large tongued", "friendly", "piano playing", "ear licking", "doe eyed", "sock sniffing", "lesbian", "hairy"}

var female = []string{"baggage", "hussy", "woman", "Duchess", "female impersonator", "nymphomaniac", "virgin", "leather freak", "home-coming queen", "defrocked nun", "bisexual budgie", "cheerleader", "office secretary", "sexual deviate", "DARPA contract monitor", "little matchgirl", "ceremonial penguin", "femme fatale", "bosses' daughter", "construction worker", "sausage abuser", "secretary", "Congressman's page", "grandmother", "penguin", "German shepherd", "stewardess", "waitress", "prostitute", "computer science group", "housewife"}

var asthe = []string{"as the"}

var madjec = []string{"thrashing", "slurping", "insatiable", "rabid", "satanic", "corpulent", "nose-grooming", "tripe-fondling", "dribbling", "spread-eagled", "orally fixated", "vile", "awesomely endowed", "handsome", "mush-brained", "tremendously hung", "three-legged", "pile-driving", "cross-dressing", "gerbil buggering", "bung-hole stuffing", "sphincter licking", "hair-pie chewing", "muff-diving", "clam shucking", "egg-sucking", "bicycle seat sniffing"}

var male = []string{"rakehell", "hunchback", "lecherous lickspittle", "archduke", "midget", "hired hand", "great Dane", "stallion", "donkey", "electric eel", "paraplegic pothead", "dirty old man", "faggot butler", "friar", "black-power advocate", "follicle fetishist", "handsome priest", "chicken flicker", "homosexual flamingo", "ex-celibate", "drug sucker", "ex-woman", "construction worker", "hair dresser", "dentist", "judge", "social worker"}

var diddled = []string{"diddled", "devoured", "fondled", "mouthed", "tongued", "lashed", "tweaked", "violated", "defiled", "irrigated", "penetrated", "ravished", "hammered", "bit", "tongue slashed", "sucked", "fucked", "rubbed", "grudge fucked", "masturbated with", "slurped"}

var her = []string{"her"}

var titadj = []string{"alabaster", "pink-tipped", "creamy", "rosebud", "moist", "throbbing", "juicy", "heaving", "straining", "mammoth", "succulent", "quivering", "rosey", "globular", "varicose", "jiggling", "bloody", "tilted", "dribbling", "oozing", "firm", "pendulous", "muscular", "bovine"}

var knockers = []string{"globes", "melons", "mounds", "buds", "paps", "chubbies", "protuberances", "treasures", "buns", "bung", "vestibule", "armpits", "tits", "knockers", "elbows", "eyes", "hooters", "jugs", "lungs", "headlights", "disk drives", "bumpers", "knees", "fried eggs", "buttocks", "charlies", "ear lobes", "bazooms", "mammaries"}

var andd = []string{"and"}

var thrust = []string{"plunged", "thrust", "squeezed", "pounded", "drove", "eased", "slid", "hammered", "squished", "crammed", "slammed", "reamed", "rammed", "dipped", "inserted", "plugged", "augured", "pushed", "ripped", "forced", "wrenched"}

var his = []string{"his"}

var dongadj = []string{"bursting", "jutting", "glistening", "Brobdingnagian", "prodigious", "purple", "searing", "swollen", "rigid", "rampaging", "warty", "steaming", "gorged", "trunklike", "foaming", "spouting", "swinish", "prosthetic", "blue veined", "engorged", "horse like", "throbbing", "humongous", "hole splitting", "serpentine", "curved", "steel encased", "glass encrusted", "knobby", "surgically altered", "metal tipped", "open sored", "rapidly dwindling", "swelling", "miniscule", "boney"}

var dong = []string{"intruder", "prong", "stump", "member", "meat loaf", "majesty", "bowsprit", "earthmover", "jackhammer", "ramrod", "cod", "jabber", "gusher", "poker", "engine", "brownie", "joy stick", "plunger", "piston", "tool", "manhood", "lollipop", "kidney prodder", "candlestick", "John Thomas", "arm", "testicles", "balls", "finger", "foot", "tongue", "dick", "one-eyed wonder worm", "canyon yodeler", "middle leg", "neck wrapper", "stick shift", "dong", "Linda Lovelace choker"}

var intoher = []string{"into her"}

var twatadj = []string{"pulsing", "hungry", "hymeneal", "palpitating", "gaping", "slavering", "welcoming", "glutted", "gobbling", "cobwebby", "ravenous", "slurping", "glistening", "dripping", "scabiferous", "porous", "soft-spoken", "pink", "dusty", "tight", "odiferous", "moist", "loose", "scarred", "weapon-less", "banana stuffed", "tire tracked", "mouse nibbled", "tightly tensed", "oft traveled", "grateful", "festering"}

var twat = []string{"swamp.", "honeypot.", "jam jar.", "butterbox.", "furburger.", "cherry pie.", "cush.", "slot.", "slit.", "cockpit.", "damp.", "furrow.", "sanctum sanctorum.", "bearded clam.", "continental divide.", "paradise valley.", "red river valley.", "slot machine.", "quim.", "palace.", "ass.", "rose bud.", "throat.", "eye socket.", "tenderness.", "inner ear.", "orifice.", "appendix scar.", "wound.", "navel.", "mouth.", "nose.", "cunt."}

func altogethernow() string {
	return choice(faster) + " " + choice(said) + " " + choice(the) + " " + choice(fadj) + " " + choice(female) + " " + choice(asthe) + " " + choice(madjec) + " " + choice(male) + " " + choice(diddled) + " " + choice(her) + " " + choice(titadj) + " " + choice(knockers) + " " + choice(andd) + " " + choice(thrust) + " " + choice(his) + " " + choice(dongadj) + " " + choice(dong) + " " + choice(intoher) + " " + choice(twatadj) + " " + choice(twat)
}

func choice(s []string) string {
	rand.Seed(time.Now().UnixNano())
	if len(s) == 1 {
		return s[0]
	}
	return s[rand.Intn(len(s)-1)]
}
