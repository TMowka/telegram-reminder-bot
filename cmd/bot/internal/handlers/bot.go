package handlers

import (
	"context"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"

	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"

	"github.com/tmowka/telegram-reminder-bot/internal/config"
)

type Bot struct {
	db     *sqlx.DB
	chatId string
}

func (b *Bot) Start(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.Start")
	defer span.End()

	if err := config.Save(ctx, b.db, config.BotStarted, true); err != nil {
		log.Println("handlers.Bot.Start : error :", err)
	}
}

func (b *Bot) Stop(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.Stop")
	defer span.End()

	if err := config.Save(ctx, b.db, config.BotStarted, false); err != nil {
		log.Println("handlers.Bot.Start : error :", err)
	}
}
