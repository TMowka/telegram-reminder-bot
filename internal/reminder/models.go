package reminder

import "time"

type Reminder struct {
	WeekdaysToSkip map[time.Weekday]struct{}
	Interval       time.Duration
	Location       *time.Location
	RemindTime     time.Time
	Started        bool

	ticker     *time.Ticker
	clearChan  chan struct{}
	remindChan chan struct{}
}
