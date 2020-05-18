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

const DATE_TIME_LAYOUT = "2 Jan 2006 15:04"

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
		err = errors.Wrap(err, "error sending telebot message")
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

	rtCh := make(chan string)
	wdtsCh := make(chan string)

	go func() {
		defer close(rtCh)

		remindTime, err := config.GetByName(ctx, b.db, config.RemindTime)
		if err != nil {
			err = errors.Wrap(err, "error getting remind time")
			log.Println("handlers.Bot.Start : error :", err)
		}

		rtCh <- remindTime.(string)
	}()

	go func() {
		defer close(wdtsCh)

		weekdaysToSkip, err := config.GetByName(ctx, b.db, config.WeekdaysToSkip)
		if err != nil {
			err = errors.Wrap(err, "error getting weekdays to skip")
			log.Println("handlers.Bot.Start : error :", err)
		}

		wdtsCh <- weekdaysToSkip.(string)
	}()

	remindTime, weekdaysToSkip := <-rtCh, <-wdtsCh

	if err := b.reminder.SetWeekdaysToSkip(weekdaysToSkip); err != nil {
		err = errors.Wrap(err, "error setting weekdays to skip")
		log.Println("handlers.Bot.Start : error :", err)
	}

	remindChan, err := b.reminder.Start(remindTime)
	if err != nil {
		err = errors.Wrap(err, "error starting reminder")
		log.Println("handlers.Bot.Start : error :", err)
		return
	}

	go func() {
		for range remindChan {
			b.notify(ctx)
		}
	}()

	if err := config.Save(ctx, b.db, config.BotStarted, true, time.Now()); err != nil {
		log.Println("handlers.Bot.Start : error :", err)
		return
	}
}

func (b *Bot) Stop(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.Stop")
	defer span.End()

	if err := b.reminder.Stop(); err != nil {
		err = errors.Wrap(err, "error stopping reminder")
		log.Println("handlers.Bot.Stop : error :", err)
		return
	}

	if err := config.Save(ctx, b.db, config.BotStarted, false, time.Now()); err != nil {
		err = errors.Wrap(err, "error saving config")
		log.Println("handlers.Bot.Stop : error :", err)
		return
	}
}

func (b *Bot) AddParticipant(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.AddParticipant")
	defer span.End()

	p := participant.NewParticipant{
		Name: strings.TrimSpace(m.Payload),
	}

	if _, err := participant.CreateOrUpdate(ctx, b.db, p, time.Now()); err != nil {
		err = errors.Wrap(err, "error adding participants")
		log.Println("handlers.Bot.AddParticipant : error :", err)
		return
	}
}

func (b *Bot) RemoveParticipant(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.RemoveParticipant")
	defer span.End()

	if err := participant.DeleteByName(ctx, b.db, m.Payload); err != nil {
		err = errors.Wrap(err, "error removing participants")
		log.Println("handlers.Bot.RemoveParticipant : error :", err)
		return
	}
}

func (b *Bot) SetRemindTime(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.SetRemindTime")
	defer span.End()

	if err := config.Save(ctx, b.db, config.RemindTime, m.Payload, time.Now()); err != nil {
		err = errors.Wrap(err, "error saving remind time")
		log.Println("handlers.Bot.SetRemindTime : error :", err)
		return
	}
}

func (b *Bot) SetRemindMessage(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.SetRemindMessage")
	defer span.End()

	if err := config.Save(ctx, b.db, config.RemindMessage, m.Payload, time.Now()); err != nil {
		err = errors.Wrap(err, "error saving remind message")
		log.Println("handlers.Bot.SetRemindMessage : error :", err)
		return
	}
}

func (b *Bot) SetWeekdaysToSkip(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.SetWeekdaysToSkip")
	defer span.End()

	if err := config.Save(ctx, b.db, config.WeekdaysToSkip, m.Payload, time.Now()); err != nil {
		err = errors.Wrap(err, "error saving weekdays to skip")
		log.Println("handlers.Bot.SetWeekdaysToSkip : error :", err)
		return
	}

	if err := b.reminder.SetWeekdaysToSkip(m.Payload); err != nil {
		err = errors.Wrap(err, "error setting weekdays to skip")
		log.Println("handlers.Bot.SetWeekdaysToSkip : error :", err)
		return
	}
}

func (b *Bot) Info(m *tb.Message) {
	ctx, span := trace.StartSpan(context.Background(), "handlers.Bot.Info")
	defer span.End()

	var weekdaysToSkip []string
	for key := range b.reminder.WeekdaysToSkip {
		weekdaysToSkip = append(weekdaysToSkip, key.String())
	}

	participantList, err := participant.List(ctx, b.db)
	if err != nil {
		err = errors.Wrap(err, "error getting participant list")
		log.Println("handlers.Bot.Info : error :", err)
		participantList = []participant.Participant{}
	}

	participants := make([]string, len(participantList))
	for i, p := range participantList {
		participants[i] = p.Name
	}

	msg := fmt.Sprintf(`
Server time: %s
Remind time: %s
Weekdays to skip: %s
Participants: %s
Reminder started: %v
`,
		time.Now().In(b.reminder.Location).Format(DATE_TIME_LAYOUT),
		b.reminder.RemindTime.In(b.reminder.Location).Format(DATE_TIME_LAYOUT),
		strings.Join(weekdaysToSkip, ", "),
		strings.Join(participants, ", "),
		b.reminder.Started,
	)

	b.send(msg)
}

func (b *Bot) Help(m *tb.Message) {
	_, span := trace.StartSpan(context.Background(), "handlers.Bot.Help")
	defer span.End()

	msg := `
/hello - Bot lifecheck
/help - Print list of available commands
/start - Start reminder with pre-configured remind time
/stop - Stop reminder
/addparticipant - Add participant to remind
/removeparticipant - Remove participant
/setremindtime - Set time of the next remind in format "HH:MM" (default interval is 24h)
/setremindmessage - Set remind message
/setweekdaystoskip - Set weekdays to skip in format "0,1,2" (0 - is Sunday)
/info - Print bot configuration and state
`

	b.send(msg)
}
