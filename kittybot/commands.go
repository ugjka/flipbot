package kitty

import (
	"github.com/bwmarrin/discordgo"
	"mvdan.cc/xurls/v2"
)

// Reply sends a message to where the message came from (user or channel)
func (bot *Bot) Reply(m *Message, text string) {
	urls := xurls.Strict().FindStringSubmatch(text)
	if len(urls) == 0 {
		m.Session.ChannelMessageSend(m.To, text)
	}
	m.Session.ChannelMessageSendEmbed(m.To, &discordgo.MessageEmbed{
		URL:         urls[0],
		Description: text,
	})
}
