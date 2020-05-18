package config

import "time"

type Name string

const (
	BotStarted     Name = "BotStarted"
	RemindTime     Name = "RemindTime"
	RemindMessage  Name = "RemindMessage"
	WeekdaysToSkip Name = "WeekdaysToSkip"
)

type config struct {
	ID        string    `db:"config_id" json:"id"`          // Unique identifier.
	Name      Name      `db:"name" json:"name"`             // Name of the config.
	CreatedAt time.Time `db:"created_at" json:"created_at"` // When the config was added.
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"` // When the config record was last modified.
}

type BooleanConfig struct {
	config
	Value bool `db:"value" json:"value"` // Boolean value of the config.
}

type IntegerConfig struct {
	config
	Value int `db:"value" json:"value"` // Integer value of the config.
}

type StringConfig struct {
	config
	Value string `db:"value" json:"value"` // String value of the config.
}
