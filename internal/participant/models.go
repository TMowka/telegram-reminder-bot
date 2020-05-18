package participant

import "time"

type Participant struct {
	ID        string    `db:"participant_id" json:"id"`
	Name      string    `db:"name" json:"name"`
	AddedAt   time.Time `db:"added_at" json:"added_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type NewParticipant struct {
	Name string `json:"name" validate:"required"`
}
