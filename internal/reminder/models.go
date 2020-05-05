package reminder

import "time"

type Reminder struct {
	RemindAt      time.Time
	RemindMessage string
	Interval      time.Duration
	Started       bool
}
