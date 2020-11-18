// Package kitty is IRCv3 enabled framework for writing IRC bots
package kitty

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	log "gopkg.in/inconshreveable/log15.v2"
	logext "gopkg.in/inconshreveable/log15.v2/ext"
)

// Bot implements an irc bot to be connected to a given server
type Bot struct {
	token    string
	outgoing chan string
	handlers []Handler
	session  *discordgo.Session
	// When did we start? Used for uptime
	started time.Time
	// Log15 loggger
	log.Logger
}

// NewBot creates a new instance of Bot
func NewBot(token string) *Bot {
	// Defaults are set here
	bot := Bot{
		started:  time.Now(),
		token:    token,
		outgoing: make(chan string, 1),
	}
	// Discard logs by default
	bot.Logger = log.New("id", logext.RandId(8))

	bot.Logger.SetHandler(log.DiscardHandler())
	return &bot
}

func (bot *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	bot.session = s
	msg := parseMessage(s, m)
	for _, handler := range bot.handlers {
		handler.Handle(bot, msg)
	}
}

// Run runs
func (bot *Bot) Run() {
	dg, err := discordgo.New("Bot " + bot.token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	dg.AddHandler(bot.messageCreate)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-sig
	dg.Close()
	bot.Info("disconnected")
	return
}

// Handler is used to subscribe and react to events on the bot Server
type Handler interface {
	Handle(*Bot, *Message)
}

// Trigger is a Handler which is guarded by a condition.
// DO NOT alter *Message in your triggers or you'll have strange things happen.
type Trigger struct {
	// Returns true if this trigger applies to the passed in message
	Condition func(*Bot, *Message) bool

	// The action to perform if Condition is true
	Action func(*Bot, *Message)
}

// AddTrigger adds a trigger to the bot's handlers
func (bot *Bot) AddTrigger(h Handler) {
	bot.handlers = append(bot.handlers, h)
}

// Handle executes the trigger action if the condition is satisfied
func (t Trigger) Handle(bot *Bot, m *Message) {
	if t.Condition(bot, m) {
		t.Action(bot, m)
	}
}

// Message represents a message received from the server
type Message struct {
	// Content generally refers to the text of a PRIVMSG
	Content string
	Command string
	Name    string
	To      string
	Session *discordgo.Session
}

// parseMessage takes a string and attempts to create a Message struct.
// Returns nil if the Message is invalid.
func parseMessage(s *discordgo.Session, m *discordgo.MessageCreate) (msg *Message) {
	msg = new(Message)
	msg.Content = m.Content
	msg.Command = "PRIVMSG"
	msg.Name = m.Author.Username
	if m.Member.Nick != "" {
		msg.Name = m.Member.Nick
	}
	msg.To = m.ChannelID
	msg.Session = s
	return msg
}
