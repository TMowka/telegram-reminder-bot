package handlers

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"github.com/tmowka/telegram-reminder-bot/internal/config"
)

type bot struct {
	db *sqlx.DB
}

func NewBot(db *sqlx.DB) *bot {
	return &bot{
		db: db,
	}
}

func (b *bot) Start(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Bot.Start")
	defer span.End()

	if err := config.Save(ctx, b.db, config.BotStarted, true); err != nil {
		return errors.Wrapf(err, "saving config: %s, %v", config.BotStarted, true)
	}

	return nil
}

func (b *bot) Stop(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Bot.Stop")
	defer span.End()

	if err := config.Save(ctx, b.db, config.BotStarted, false); err != nil {
		return errors.Wrapf(err, "saving config: %s, %v", config.BotStarted, false)
	}

	return nil
}
