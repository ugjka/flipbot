package kitty

import "github.com/bwmarrin/discordgo"

// Reply sends a message to where the message came from (user or channel)
func (bot *Bot) Reply(m *Message, text string) {
	if Self(m) {
		return
	}
	m.Session.ChannelMessageSend(m.To, text)
}

// Rich is rich
type Rich struct {
	URL         string
	Title       string
	Description string
	IconURL     string
}

// ReplyRich sends a message to where the message came from (user or channel)
func (bot *Bot) ReplyRich(m *Message, r Rich) {
	if Self(m) {
		return
	}
	m.Session.ChannelMessageSendEmbed(m.To, &discordgo.MessageEmbed{
		URL:         r.URL,
		Title:       r.Title,
		Description: r.Description,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: r.IconURL,
		},
	})
}

// Self checks for itself
func Self(m *Message) bool {
	return m.m.Author.ID == m.Session.State.User.ID
}
