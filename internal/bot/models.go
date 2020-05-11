package bot

import (
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/tmowka/telegram-reminder-bot/internal/reminder"
)

type Bot struct {
	Chats []tb.Recipient

	participants Participants
	reminder     *reminder.Reminder
	location     *time.Location
}

type Participants map[string]time.Time

func (p Participants) Print() string {
	var participants []string
	for key := range p {
		participants = append(participants, key)
	}
	return strings.Join(participants, ", ")
}

type Chat struct {
	chatId string
}

func (c *Chat) Recipient() string {
	return c.chatId
}
