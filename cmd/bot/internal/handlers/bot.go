package handlers

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/tmowka/telegram-reminder-bot/internal/config"
	"github.com/tmowka/telegram-reminder-bot/internal/participant"
	"github.com/tmowka/telegram-reminder-bot/internal/reminder"
)

type Bot struct {
	db       *sqlx.DB
	chat     *chat
	telebot  *tb.Bot
	reminder *reminder.Reminder
}

type chat struct {
	id string
}

func (c *chat) Recipient() string {
	return c.id
}

func (b *Bot) send(msg string) {
	if _, err := b.telebot.Send(b.chat, msg); err != nil {
		log.Println("handlers.Bot.send : error :", err)
	}
}

func (b *Bot) notify() {
	b.send("remind message")
}

func (b *Bot) Hello(m *tb.Message) {
	_, span := trace.StartSpan(context.Background(), "handlers.Bot.Hello")
	defer span.End()

	b.send("Hello World!")
}

func (b *Bot) Start(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.Start")
	defer span.End()

	if err := config.Save(ctx, b.db, config.BotStarted, true, time.Now()); err != nil {
		log.Println("handlers.Bot.Start : error :", err)
	}

	remindChan, err := b.reminder.Start(time.Now().Add(3 * time.Second))
	if err != nil {
		err = errors.Wrap(err, "error starting reminder")
		log.Println("handlers.Bot.Start : error :", err)
	}

	go func() {
		for _ = range remindChan {
			b.notify()
		}
	}()
}

func (b *Bot) Stop(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.Stop")
	defer span.End()

	if err := config.Save(ctx, b.db, config.BotStarted, false, time.Now()); err != nil {
		err = errors.Wrap(err, "error saving config")
		log.Println("handlers.Bot.Stop : error :", err)
	}

	if err := b.reminder.Stop(); err != nil {
		err = errors.Wrap(err, "error stopping reminder")
		log.Println("handlers.Bot.Stop : error :", err)
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

func (b *Bot) SetRemindMessage(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.SetRemindMessage")
	defer span.End()

	if err := config.Save(ctx, b.db, config.RemindMessage, m.Payload, time.Now()); err != nil {
		log.Println("handlers.Bot.SetRemindMessage : error :", err)
	}
}

func (b *Bot) SetWeekdaysToSkip(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.SetWeekdaysToSkip")
	defer span.End()

	if err := config.Save(ctx, b.db, config.WeekdaysToSkip, m.Payload, time.Now()); err != nil {
		log.Println("handlers.Bot.SetWeekdaysToSkip : error :", err)
	}
}
