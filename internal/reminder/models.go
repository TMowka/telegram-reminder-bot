package reminder

import "time"

type Reminder struct {
	RemindTime    time.Time
	RemindMessage string
	RemindChan    chan string
	Started       bool

	weekdaysToSkip map[time.Weekday]struct{}
	interval       time.Duration
	ticker         *time.Ticker
	clearChan      chan bool
}
