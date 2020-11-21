package kitty

// Reply sends a message to where the message came from (user or channel)
func (bot *Bot) Reply(m *Message, text string) {
	bot.Info("Discord", " chan", m.To)
	m.Session.ChannelMessageSend(m.To, text)
}
