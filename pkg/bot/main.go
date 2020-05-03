package bot

import (
	"time"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Runnable interface {
	Run() error
}

type bot struct {
	url     string
	token   string
	timeout time.Duration
}

func New(url string, token string, timeout time.Duration) Runnable {
	return &bot{url: url, token: token, timeout: timeout}
}

func (b *bot) Run() error {
	bSettings := tb.Settings{
		URL:    b.url,
		Token:  b.token,
		Poller: &tb.LongPoller{Timeout: b.timeout},
	}

	bot, err := tb.NewBot(bSettings)
	if err != nil {
		return errors.Wrapf(err, "Run->NewBot(%+v)", bSettings)
	}

	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "Hello World!")
	})

	bot.Start()

	return nil
}
