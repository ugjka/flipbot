package kitty

// Msg sends a message to 'who' (user or channel)
func (bot *Bot) Msg(who, text string) {
	bot.Send(text)
}

// Reply sends a message to where the message came from (user or channel)
func (bot *Bot) Reply(m *Message, text string) {
	_ = m
	bot.Send(text)
}

// Send any command to the server
func (bot *Bot) Send(command string) {
	bot.outgoing <- &Message{
		Command: "PRIVMSG",
		Content: command}
}
