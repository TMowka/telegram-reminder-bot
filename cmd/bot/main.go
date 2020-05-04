package main

import (
	"flag"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/tmowka/telegram-reminder-bot/internal/bot"
)

type config struct {
	url     string
	token   string
	chatIds string
	timeout uint64
}

func main() {
	var cfg config
	flag.StringVar(&cfg.token, "token", "", "telegram bot token")
	flag.StringVar(&cfg.chatIds, "chat-ids", "", "telegram chat ids list")
	flag.Parse()

	telebot, err := tb.NewBot(tb.Settings{
		Token:  cfg.token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	b := bot.NewBot(strings.Split(cfg.chatIds, ","))

	b.Run(telebot)
}
