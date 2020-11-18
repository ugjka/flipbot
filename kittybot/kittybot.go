// Package kitty is IRCv3 enabled framework for writing IRC bots
package kitty

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	log "gopkg.in/inconshreveable/log15.v2"
	logext "gopkg.in/inconshreveable/log15.v2/ext"
)

// Bot implements an irc bot to be connected to a given server
type Bot struct {
	token       string
	outgoing    chan *Message
	incoming    chan *discordgo.MessageCreate
	handlers    []Handler
	session     *discordgo.Session
	sessionOnce sync.Once
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
		outgoing: make(chan *Message),
		incoming: make(chan *discordgo.MessageCreate),
	}
	// Discard logs by default
	bot.Logger = log.New("id", logext.RandId(8))

	bot.Logger.SetHandler(log.DiscardHandler())
	return &bot
}

// Handles message speed throtling
func (bot *Bot) handleOutgoingMessages() {
	for {
		msg := <-bot.outgoing
		bot.session.ChannelMessageSend(msg.To, msg.Content)
	}
}

func (bot *Bot) handleIncomingMessages() {
	for {
		msg := parseMessage(<-bot.incoming)
		go func() {
			for _, h := range bot.handlers {
				go h.Handle(bot, msg)
			}
		}()
	}
}

func (bot *Bot) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	bot.sessionOnce.Do(func() {
		bot.session = s
	})
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	bot.incoming <- m

	// If the message is "pong" reply with "Ping!"
	msg := <-bot.outgoing
	s.ChannelMessageSend(msg.To, msg.Content)
}

// Run runs
func (bot *Bot) Run() {
	dg, err := discordgo.New("Bot " + bot.token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)
	bot.Debug("starting bot goroutines")
	go bot.handleIncomingMessages()
	go bot.handleOutgoingMessages()

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-sig
	bot.Info("disconnected")
	return
}

// Uptime returns the uptime of the bot
func (bot *Bot) Uptime() string {
	return fmt.Sprintf("Started: %s, Uptime: %s", bot.started, time.Since(bot.started))
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
}

// parseMessage takes a string and attempts to create a Message struct.
// Returns nil if the Message is invalid.
func parseMessage(raw *discordgo.MessageCreate) (m *Message) {
	m = new(Message)
	m.Content = raw.Content
	m.Command = "PRIVMSG"
	m.Name = raw.Author.Username
	m.To = raw.ChannelID
	return m
}