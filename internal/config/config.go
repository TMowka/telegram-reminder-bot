package config

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific Config is requested but does not exist.
	ErrNotFound = errors.New("Config not found")
)

func GetByName(ctx context.Context, db *sqlx.DB, name Name) (interface{}, error) {
	ctx, span := trace.StartSpan(ctx, "internal.config.GetByName")
	defer span.End()

	var bc BooleanConfig
	//var ic IntegerConfig
	//var sc StringConfig

	const q = `select * from config
		where name = $1`

	switch name {
	case BotStarted:
		if err := db.GetContext(ctx, &bc, q, name); err != nil {
			if err == sql.ErrNoRows {
				return nil, ErrNotFound
			}

			return nil, errors.Wrap(err, "selecting boolean config by name")
		}
		return &bc, nil
	default:
		return nil, errors.New("unknown config name to select")
	}
}

func Save(ctx context.Context, db *sqlx.DB, name Name, val interface{}, now time.Time) error {
	ctx, span := trace.StartSpan(ctx, "internal.config.Save")
	defer span.End()

	const updateQ = `update config
		set value = $1, updated_at = $2
		where name = $3`
	const insertQ = `insert into config
		(config_id, name, value, created_at, updated_at)
		values ($1, $2, $3, $4, $5)`

	switch name {
	case BotStarted:
		cfg := &BooleanConfig{
			config: config{
				ID:   uuid.New().String(),
				Name: name,
			},
			Value: val.(bool),
		}

		res, err := db.ExecContext(ctx, updateQ,
			cfg.Value, now.UTC(), cfg.Name,
		)
		if err != nil {
			return errors.Wrap(err, "updating config")
		}

		upd, err := res.RowsAffected()
		if upd > 0 {
			return nil
		}

		_, err = db.ExecContext(ctx, insertQ,
			cfg.ID, cfg.Name, cfg.Value, now.UTC(), now.UTC())
		if err != nil {
			return errors.Wrap(err, "inserting config")
		}
		return nil
	default:
		return errors.New("unknown config name to insert")
	}
}
