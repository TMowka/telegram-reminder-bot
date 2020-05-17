package handlers

import (
	"time"

	"github.com/jmoiron/sqlx"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/tmowka/telegram-reminder-bot/internal/reminder"
)

func Telebot(db *sqlx.DB, telebot *tb.Bot, chatId string) error {
	r := reminder.New(5 * time.Second)

	b := Bot{
		db: db,
		chat: &chat{
			id: chatId,
		},
		telebot:  telebot,
		reminder: r,
	}

	telebot.Handle("/hello", b.Hello)
	telebot.Handle("/start", b.Start)
	telebot.Handle("/stop", b.Stop)
	telebot.Handle("/addparticipant", b.AddParticipant)
	telebot.Handle("/removeparticipant", b.RemoveParticipant)
	telebot.Handle("/setremindtime", b.SetRemindTime)
	telebot.Handle("/setremindmessage", b.SetRemindMessage)
	telebot.Handle("/setweekdaystoskip", b.SetWeekdaysToSkip)

	return nil
}
