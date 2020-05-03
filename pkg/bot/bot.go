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
	Run() error
}

type bot struct {
	telebot      *tb.Bot
	chat         tb.Recipient
	participants map[string]participant
	remindAt     time.Time
	interval     time.Duration
}

var ticker *time.Ticker

func New(token string, chatId string) (Runnable, error) {
	bSettings := tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}

	telebot, err := tb.NewBot(bSettings)

	if err != nil {
		return nil, errors.Wrapf(err, "New->NewBot(%+v)", bSettings)
	}

	return &bot{
		chat: &chat{
			chatId: chatId,
		},
		telebot:      telebot,
		participants: make(map[string]participant),
		interval:     24 * time.Hour,
	}, nil
}

func (b *bot) remind() {
	b.remindAt = time.Now().Add(b.interval)
	b.send(fmt.Sprintf("Reminder: %v\nFill in the project server, please.", b.printParticipants()))
}

func (b *bot) start() {
	if ticker != nil {
		ticker.Stop()
	}

	ticker = time.NewTicker(time.Second)

	for now := range ticker.C {
		if now.Unix() >= b.remindAt.Unix() {
			b.remind()
		}
	}
}

func (b *bot) stop() {
	if ticker != nil {
		ticker.Stop()
		ticker = nil
	}
}

func (b *bot) send(msg string) {
	_, _ = b.telebot.Send(b.chat, msg)
}

func (b *bot) addHandler(endpoint string, handler func(m *tb.Message)) {
	b.telebot.Handle(endpoint, func(m *tb.Message) {
		if !m.Private() {
			return
		}

		handler(m)
	})
}

func (b *bot) Run() error {
	b.addHandler("/hello", func(m *tb.Message) {
		b.send("Hello World!")
	})

	b.addHandler("/add_participant", func(m *tb.Message) {
		b.addParticipant(m.Payload)
		b.send("Participants: " + b.printParticipants())
	})

	b.addHandler("/add_participants", func(m *tb.Message) {
		participants := strings.Split(m.Payload, ", ")
		for _, p := range participants {
			b.addParticipant(p)
		}
		b.send("Participants: " + b.printParticipants())
	})

	b.addHandler("/remove_participant", func(m *tb.Message) {
		b.removeParticipant(m.Payload)
		b.send("Participants: " + b.printParticipants())
	})

	b.addHandler("/set_remind_time", func(m *tb.Message) {
		remindAtArr := strings.Split(m.Payload, ":")

		if len(remindAtArr) != 2 {
			return
		}

		remindAtH, err := strconv.Atoi(remindAtArr[0])
		if err != nil {
			return
		}

		remindAtM, err := strconv.Atoi(remindAtArr[1])
		if err != nil {
			return
		}

		now := time.Now()
		remindAt := time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			remindAtH,
			remindAtM,
			0,
			0,
			time.UTC,
		)

		if remindAtH < now.Hour() {
			remindAt = remindAt.Add(24 * time.Hour)
		}

		b.remindAt = remindAt

		if ticker != nil {
			b.send(fmt.Sprintf("Next remind at: %s", b.remindAt))
		}
	})

	b.addHandler("/start", func(m *tb.Message) {
		if b.remindAt.IsZero() {
			b.send("Set remind time first")
			return
		}

		b.send(fmt.Sprintf("Next remind at: %s", b.remindAt))
		b.start()
	})

	b.addHandler("/stop", func(m *tb.Message) {
		b.stop()
	})

	b.telebot.Start()

	return nil
}
