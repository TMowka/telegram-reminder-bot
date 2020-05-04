package bot

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Runnable interface {
	Run(telebot *tb.Bot)
}

func NewBot(chatIds []string) Runnable {
	chats := make([]tb.Recipient, len(chatIds))
	for i, cId := range chatIds {
		chats[i] = &chat{
			chatId: cId,
		}
	}

	return &bot{
		chats:         chats,
		participants:  make(map[string]time.Time),
		interval:      24 * time.Hour,
		remindMessage: "Fill in the project server please.  🤞",
	}
}

func (b *bot) sendMessage(telebot *tb.Bot, msg string) {
	for _, chat := range b.chats {
		if _, err := telebot.Send(chat, msg); err != nil {
			fmt.Printf("error occured while sending the message: %v",
				errors.Wrapf(err, "sendMessage->telebot.Send(%+v, %s)", chat, msg))
		}
	}
}

func (b *bot) addParticipants(participants []string) {
	for _, p := range participants {
		fmtParticipant := strings.TrimSpace(p)
		if len(fmtParticipant) > 0 {
			b.participants[fmtParticipant] = time.Now()
		}
	}
}

func (b *bot) removeParticipant(participant string) {
	delete(b.participants, participant)
}

func (b *bot) printParticipants() string {
	var participants []string
	for key, _ := range b.participants {
		participants = append(participants, key)
	}
	return strings.Join(participants, " ")
}

func (b *bot) remind(telebot *tb.Bot) {
	b.remindAt = time.Now().Add(b.interval)
	b.sendMessage(telebot, fmt.Sprintf("%v\n%s",
		b.printParticipants(), b.remindMessage))
	fmt.Printf("Next remind at: %s", b.remindAt)
}

func parseTime(raw string) time.Time {
	hmArr := strings.Split(raw, ":")

	if len(hmArr) != 2 {
		return time.Time{}
	}

	h, err := strconv.Atoi(hmArr[0])
	if err != nil {
		return time.Time{}
	}

	m, err := strconv.Atoi(hmArr[1])
	if err != nil {
		return time.Time{}
	}

	now := time.Now()
	remindAt := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		h,
		m,
		0,
		0,
		time.Local,
	)

	if remindAt.Unix() < now.Unix() {
		remindAt = remindAt.Add(24 * time.Hour)
	}

	return remindAt
}

var ticker *time.Ticker
var clearChan chan bool

func (b *bot) Run(telebot *tb.Bot) {
	telebot.Handle("/hello", func(m *tb.Message) {
		b.sendMessage(telebot, "Hello World!")
	})

	telebot.Handle("/participants", func(m *tb.Message) {
		b.sendMessage(telebot, b.printParticipants())
	})

	telebot.Handle("/add_participants", func(m *tb.Message) {
		b.addParticipants(strings.Split(m.Payload, ","))
	})

	telebot.Handle("/remove_participant", func(m *tb.Message) {
		b.removeParticipant(m.Payload)
	})

	telebot.Handle("/message", func(m *tb.Message) {
		message := strings.TrimSpace(m.Payload)
		if len(message) > 0 {
			b.remindMessage = message
		}
	})

	telebot.Handle("/start", func(m *tb.Message) {
		if len(m.Payload) > 0 {
			b.remindAt = parseTime(m.Payload)
		}

		if b.remindAt.IsZero() {
			b.sendMessage(telebot, "Remind time is not set")
			return
		}

		if ticker != nil {
			ticker.Stop()
		}

		ticker = time.NewTicker(time.Second)
		clearChan = make(chan bool)

		b.started = true

		go func() {
			for {
				select {
				case <-ticker.C:
					if time.Now().Unix() >= b.remindAt.Unix() {
						b.remind(telebot)
					}
				case <-clearChan:
					ticker.Stop()
					return
				}
			}
		}()
	})

	telebot.Handle("/stop", func(m *tb.Message) {
		clearChan <- true
		b.started = false
	})

	telebot.Handle("/info", func(m *tb.Message) {
		b.sendMessage(telebot, fmt.Sprintf("Server time: %s\nNext remind at: %s\nStarted: %v",
			time.Now(), b.remindAt, b.started))
	})

	telebot.Start()
}
