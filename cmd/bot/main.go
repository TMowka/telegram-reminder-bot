package main

import (
	"flag"
	"time"

	"github.com/tmowka/telegram-reminder-bot/pkg/bot"
)

type config struct {
	url     string
	token   string
	timeout uint64
}

func main() {
	var cfg config
	flag.StringVar(&cfg.url, "url", "", "telegram bot url")
	flag.StringVar(&cfg.token, "token", "", "telegram bot token")
	flag.Uint64Var(&cfg.timeout, "timeout", 10, "telegram bot poller timeout")
	flag.Parse()

	b := bot.New(cfg.url, cfg.token, time.Duration(cfg.timeout)*time.Second)

	if err := b.Run(); err != nil {
		panic(err)
	}
}
