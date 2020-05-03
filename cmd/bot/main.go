package main

import (
	"flag"
	"github.com/tmowka/telegram-reminder-bot/pkg/bot"
)

type config struct {
	url     string
	token   string
	chatId  string
	timeout uint64
}

func main() {
	var cfg config
	flag.StringVar(&cfg.token, "token", "", "telegram bot token")
	flag.StringVar(&cfg.chatId, "chat-id", "", "telegram chat id")
	flag.Parse()

	b, err := bot.New(cfg.token, cfg.chatId)
	if err != nil {
		panic(err)
	}

	if err := b.Run(); err != nil {
		panic(err)
	}
}
