package bot

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/tmowka/telegram-reminder-bot/internal/reminder"
)

type Bot struct {
	Chats []tb.Recipient

	participants map[string]time.Time
	reminder     *reminder.Reminder
	location     *time.Location
}

type Chat struct {
	chatId string
}

func (c *Chat) Recipient() string {
	return c.chatId
}
