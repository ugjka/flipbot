package kitty

import "github.com/bwmarrin/discordgo"

// Reply sends a message to where the message came from (user or channel)
func (bot *Bot) Reply(m *Message, text string) {
	//bot.Info("Discord", " chan", m.To)
	m.Session.ChannelMessageSend(m.To, text)
}

// ReplyMP3 sends a message to where the message came from (user or channel)
func (bot *Bot) ReplyMP3(m *Message, text string) {
	//bot.Info("Discord", " chan", m.To)
	m.Session.ChannelMessageSendEmbed(m.To, &discordgo.MessageEmbed{
		URL:  text,
		Type: discordgo.EmbedTypeVideo,
	})
}
