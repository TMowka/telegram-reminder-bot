package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/tmowka/telegram-reminder-bot/internal/reminder"
)

func NewBot(chatIds []string, reminder *reminder.Reminder) *Bot {
	chats := make([]tb.Recipient, len(chatIds))
	for i, cId := range chatIds {
		chats[i] = &Chat{
			chatId: cId,
		}
	}

	return &Bot{
		Chats:        chats,
		participants: make(map[string]time.Time),
		reminder:     reminder,
	}
}

func (b *Bot) Run(telebot *tb.Bot) {
	telebot.Handle("/hello", func(m *tb.Message) {
		b.sendMessage(telebot, "Hello World!")
	})

	telebot.Handle("/add_participants", func(m *tb.Message) {
		b.addParticipants(strings.Split(m.Payload, ","))
	})

	telebot.Handle("/remove_participant", func(m *tb.Message) {
		b.removeParticipant(m.Payload)
	})

	telebot.Handle("/set_remind_time", func(m *tb.Message) {
		hmArr := strings.Split(m.Payload, ":")

		if len(hmArr) != 2 {
			fmt.Printf("invalid remind time")
			return
		}

		hour, err := strconv.Atoi(hmArr[0])
		if err != nil {
			fmt.Printf("invalid remind \"hour\" value: %v", err)
			return
		}

		min, err := strconv.Atoi(hmArr[1])
		if err != nil {
			fmt.Printf("invalid remind \"minute\" value: %v", err)
			return
		}

		loc, _ := time.LoadLocation("Europe/Minsk")

		now := time.Now()
		remindAt := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			hour,
			min,
			0,
			0,
			loc,
		).UTC()

		if remindAt.Unix() < now.Unix() {
			remindAt = remindAt.Add(24 * time.Hour)
		}

		b.reminder.RemindAt = remindAt
	})

	telebot.Handle("/set_remind_message", func(m *tb.Message) {
		message := strings.TrimSpace(m.Payload)
		if len(message) > 0 {
			b.reminder.RemindMessage = message
		}
	})

	telebot.Handle("/start", func(m *tb.Message) {
		remindChan := make(chan string)
		if err := b.reminder.Start(remindChan); err != nil {
			fmt.Printf("error occured while starting reminder: %v", err)
		}

		for remindMsg := range remindChan {
			b.sendMessage(telebot, fmt.Sprintf("%v\n%s",
				b.printParticipants(), remindMsg))
		}
	})

	telebot.Handle("/stop", func(m *tb.Message) {
		if err := b.reminder.Stop(); err != nil {
			fmt.Printf("error occured while stoping reminder: %v", err)
		}
	})

	telebot.Handle("/info", func(m *tb.Message) {
		loc, _ := time.LoadLocation("Europe/Minsk")
		b.sendMessage(telebot, fmt.Sprintf(
			"Participants: %s\n"+
				"Server time: %s\n"+
				"Next remind at: %s\n"+
				"Remind message: %s\n"+
				"Started: %v",
			b.printParticipants(),
			time.Now().In(loc),
			b.reminder.RemindAt.In(loc),
			b.reminder.RemindMessage,
			b.reminder.Started),
		)
	})

	telebot.Start()
}

func (b *Bot) sendMessage(telebot *tb.Bot, msg string) {
	for _, chat := range b.Chats {
		if _, err := telebot.Send(chat, msg); err != nil {
			fmt.Printf("error occured while sending the message: %v",
				errors.Wrapf(err, "sendMessage->telebot.Send(%+v, %s)", chat, msg))
		}
	}
}

func (b *Bot) addParticipants(participants []string) {
	for _, p := range participants {
		fmtParticipant := strings.TrimSpace(p)
		if len(fmtParticipant) > 0 {
			b.participants[fmtParticipant] = time.Now()
		}
	}
}

func (b *Bot) removeParticipant(participant string) {
	delete(b.participants, participant)
}

func (b *Bot) printParticipants() string {
	var participants []string
	for key, _ := range b.participants {
		participants = append(participants, key)
	}
	return strings.Join(participants, ", ")
}
