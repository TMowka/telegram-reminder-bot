package bot

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Config struct {
	Token    string
	Location string
}

func Create(cfg Config) (*tb.Bot, error) {
	telebot, err := tb.NewBot(tb.Settings{
		Token:  cfg.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	return telebot, err
}
