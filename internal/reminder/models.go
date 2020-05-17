package reminder

import "time"

type Reminder struct {
	WeekdaysToSkip map[time.Weekday]struct{}
	Interval       time.Duration

	remindTime time.Time
	ticker     *time.Ticker
	clearChan  chan struct{}
	remindChan chan struct{}
	started    bool
}
