package participant

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"
)

// Predefined errors identify expected failure conditions.
var (
	// ErrInvalidID is used when an invalid UUID is provided.
	ErrInvalidID = errors.New("ID is not in its proper form")
)

func List(ctx context.Context, db *sqlx.DB) ([]Participant, error) {
	ctx, span := trace.StartSpan(ctx, "internal.participant.List")
	defer span.End()

	var participants []Participant
	const q = `select * from participants`

	if err := db.SelectContext(ctx, &participants, q); err != nil {
		return nil, errors.Wrap(err, "selecting participants")
	}

	return participants, nil
}

func CreateOrUpdate(ctx context.Context, db *sqlx.DB, participant NewParticipant, now time.Time) (*Participant, error) {
	ctx, span := trace.StartSpan(ctx, "internal.participant.Create")
	defer span.End()

	p := Participant{
		ID:        uuid.New().String(),
		Name:      participant.Name,
		AddedAt:   now.UTC(),
		UpdatedAt: now.UTC(),
	}

	const updateQ = `update participants
		set updated_at = $1
		where name = $2`
	const insertQ = `insert into participants
		(participant_id, name, added_at, updated_at)
		values ($1, $2, $3, $4)`

	res, err := db.ExecContext(ctx, updateQ,
		p.UpdatedAt, p.Name,
	)
	if err != nil {
		return nil, errors.Wrap(err, "updating participant")
	}

	upd, err := res.RowsAffected()
	if upd > 0 {
		return &p, nil
	}

	_, err = db.ExecContext(ctx, insertQ,
		p.ID, p.Name, p.AddedAt, p.UpdatedAt,
	)
	if err != nil {
		return nil, errors.Wrap(err, "inserting participant")
	}

	return &p, nil
}

func DeleteByName(ctx context.Context, db *sqlx.DB, name string) error {
	ctx, span := trace.StartSpan(ctx, "internal.participant.DeleteByName")
	defer span.End()

	const q = `delete from participants
		where name = $1`

	if _, err := db.ExecContext(ctx, q, name); err != nil {
		return errors.Wrapf(err, "deleting participant by name %s", name)
	}

	return nil
}
