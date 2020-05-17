package handlers

import (
	"context"
	"fmt"
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
		return
	}

	log.Println("handlers.Bot.send : msg :", msg)
}

func (b *Bot) notify(ctx context.Context) {
	pCh := make(chan []participant.Participant)
	rmCh := make(chan string)

	go func() {
		defer close(pCh)

		participants, err := participant.List(ctx, b.db)
		if err != nil {
			err = errors.Wrap(err, "error getting participants")
			log.Println("handlers.Bot.notify : error :", err)

			participants = []participant.Participant{}
		}

		pCh <- participants
	}()

	go func() {
		defer close(rmCh)

		remindMessage, err := config.GetByName(ctx, b.db, config.RemindMessage)
		if err != nil {
			err = errors.Wrap(err, "error getting remind message")
			log.Println("handlers.Bot.notify : error :", err)

			remindMessage = "Fill in project server, please!"
		}

		rmCh <- remindMessage.(string)
	}()

	participants, remindMessage := <-pCh, <-rmCh

	pNames := make([]string, len(participants))
	for i, p := range participants {
		pNames[i] = p.Name
	}

	var pMessage string
	if len(pNames) > 0 {
		pMessage = strings.Join(pNames, ", ")
	}

	b.send(fmt.Sprintf("%s\n%s", pMessage, remindMessage))
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
		for range remindChan {
			b.notify(ctx)
		}
	}()
}

func (b *Bot) Stop(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.Stop")
	defer span.End()

	if err := b.reminder.Stop(); err != nil {
		err = errors.Wrap(err, "error stopping reminder")
		log.Println("handlers.Bot.Stop : error :", err)
	}

	if err := config.Save(ctx, b.db, config.BotStarted, false, time.Now()); err != nil {
		err = errors.Wrap(err, "error saving config")
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
