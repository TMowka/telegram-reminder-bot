package main

import (
	"flag"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/tmowka/telegram-reminder-bot/internal/bot"
	"github.com/tmowka/telegram-reminder-bot/internal/reminder"
)

type config struct {
	token    string
	chatIds  string
	location string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.token, "token", "", "telegram bot token")
	flag.StringVar(&cfg.chatIds, "chat-ids", "", "telegram chat ids list")
	flag.StringVar(&cfg.location, "location", "", "time zone location")
	flag.Parse()

	telebot, err := tb.NewBot(tb.Settings{
		Token:  cfg.token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		panic(err)
	}

	rmdChan := make(chan string)
	rmd := reminder.NewReminder("Fill in the project server please!", rmdChan)

	loc, err := time.LoadLocation(cfg.location)
	if err != nil {
		panic(err)
	}

	b := bot.NewBot(strings.Split(cfg.chatIds, ","), loc, rmd)

	b.Run(telebot)
}
