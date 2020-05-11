package config

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrNotFound is used when a specific Config is requested but does not exist.
	ErrNotFound = errors.New("Config not found")

	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper form")

	// ErrForbidden occurs when a user tries to do something that is forbidden to
	// them according to our access control policies.
	ErrForbidden = errors.New("Attempted action is not allowed")
)

func GetByName(ctx context.Context, db *sqlx.DB, name Name) (interface{}, error) {
	ctx, span := trace.StartSpan(ctx, "internal.config.GetByName")
	defer span.End()

	var bc BooleanConfig
	//var ic IntegerConfig
	//var sc StringConfig

	const q = `select c.* from config as c
		where c.name = $1`

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

func Save(ctx context.Context, db *sqlx.DB, name Name, val interface{}) error {
	ctx, span := trace.StartSpan(ctx, "internal.config.Save")
	defer span.End()

	const updateQ = `update config
		set value = $1
		where name = $2`
	const insertQ = `insert into config
		(config_id, name, value)
		values ($1, $2, $3)`

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
			cfg.Value, cfg.Name,
		)
		if err != nil {
			return errors.Wrap(err, "updating config")
		}

		upd, err := res.RowsAffected()
		if upd > 0 {
			return nil
		}

		_, err = db.ExecContext(ctx, insertQ,
			cfg.ID, cfg.Name, cfg.Value)
		if err != nil {
			return errors.Wrap(err, "inserting config")
		}
		return nil
	default:
		return errors.New("unknown config name to insert")
	}
}
