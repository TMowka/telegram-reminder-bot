package bot

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type bot struct {
	chats         []tb.Recipient
	participants  map[string]time.Time
	remindAt      time.Time
	remindMessage string
	interval      time.Duration
	started       bool
}

type chat struct {
	chatId string
}

func (c *chat) Recipient() string {
	return c.chatId
}
