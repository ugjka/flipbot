package kitty

import "github.com/bwmarrin/discordgo"

// Reply sends a message to where the message came from (user or channel)
func (bot *Bot) Reply(m *Message, text string) {
	//bot.Info("Discord", " chan", m.To)
	m.Session.ChannelMessageSend(m.To, text)
}

// Rich is rich
type Rich struct {
	URL         string
	Title       string
	Description string
}

// ReplyRich sends a message to where the message came from (user or channel)
func (bot *Bot) ReplyRich(m *Message, r Rich) {
	//bot.Info("Discord", " chan", m.To)
	m.Session.ChannelMessageSendEmbed(m.To, &discordgo.MessageEmbed{
		URL:         r.URL,
		Title:       r.Title,
		Description: r.Description,
		Type:        discordgo.EmbedTypeVideo,
	})
}
