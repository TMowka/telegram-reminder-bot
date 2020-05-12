package handlers

import (
	"github.com/jmoiron/sqlx"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Telebot(db *sqlx.DB, telebot *tb.Bot, chatId string) error {
	b := Bot{db: db, chatId: chatId}
	telebot.Handle("/start", b.Start)
	telebot.Handle("/stop", b.Stop)

	return nil
}
