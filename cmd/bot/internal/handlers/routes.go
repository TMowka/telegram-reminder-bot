package handlers

import (
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/tmowka/telegram-reminder-bot/internal/reminder"
)

func Telebot(db *sqlx.DB, telebot *tb.Bot, chatId string, location string) error {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return errors.Wrap(err, "error loading location")
	}

	r := reminder.New(24*time.Hour, loc)

	b := Bot{
		db: db,
		chat: &chat{
			id: chatId,
		},
		telebot:  telebot,
		reminder: r,
	}

	telebot.Handle("/hello", b.Hello)
	telebot.Handle("/help", b.Help)
	telebot.Handle("/start", b.Start)
	telebot.Handle("/stop", b.Stop)
	telebot.Handle("/addparticipant", b.AddParticipant)
	telebot.Handle("/removeparticipant", b.RemoveParticipant)
	telebot.Handle("/setremindtime", b.SetRemindTime)
	telebot.Handle("/setremindmessage", b.SetRemindMessage)
	telebot.Handle("/setweekdaystoskip", b.SetWeekdaysToSkip)
	telebot.Handle("/info", b.Info)

	return nil
}
