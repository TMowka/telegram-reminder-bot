package bot

type chat struct {
	chatId string
}

func (c *chat) Recipient() string {
	return c.chatId
}
