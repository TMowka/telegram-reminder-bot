package handlers

import (
	"context"
	"github.com/tmowka/telegram-reminder-bot/internal/participant"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/tmowka/telegram-reminder-bot/internal/config"
)

type Bot struct {
	db     *sqlx.DB
	chatId string
}

func (b *Bot) Start(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.Start")
	defer span.End()

	if err := config.Save(ctx, b.db, config.BotStarted, true, time.Now()); err != nil {
		log.Println("handlers.Bot.Start : error :", err)
	}
}

func (b *Bot) Stop(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.Stop")
	defer span.End()

	if err := config.Save(ctx, b.db, config.BotStarted, false, time.Now()); err != nil {
		log.Println("handlers.Bot.Start : error :", err)
	}
}

func (b *Bot) AddParticipant(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.AddParticipant")
	defer span.End()

	p := participant.NewParticipant{
		Name: strings.TrimSpace(m.Payload),
	}

	if _, err := participant.CreateOrUpdate(ctx, b.db, p, time.Now()); err != nil {
		log.Println("handlers.Bot.AddParticipant : error :", err)
	}
}

func (b *Bot) RemoveParticipant(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.RemoveParticipant")
	defer span.End()

	if err := participant.DeleteByName(ctx, b.db, m.Payload); err != nil {
		log.Println("handlers.Bot.RemoveParticipant : error :", err)
	}
}

func (b *Bot) SetRemindTime(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.SetRemindTime")
	defer span.End()

	if err := config.Save(ctx, b.db, config.RemindTime, m.Payload, time.Now()); err != nil {
		log.Println("handlers.Bot.SetRemindTime : error :", err)
	}
}

func (b *Bot) SetWeekdaysToSkip(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.SetWeekdaysToSkip")
	defer span.End()

	if err := config.Save(ctx, b.db, config.RemindTime, m.Payload, time.Now()); err != nil {
		log.Println("handlers.Bot.SetWeekdaysToSkip : error :", err)
	}
}
