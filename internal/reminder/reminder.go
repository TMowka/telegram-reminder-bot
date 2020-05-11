package reminder

import (
	"strings"
	"time"

	"github.com/pkg/errors"
)

func NewReminder(remindMessage string, remindChan chan string) *Reminder {
	return &Reminder{
		RemindMessage: remindMessage,
		RemindChan:    remindChan,

		weekdaysToSkip: make(map[time.Weekday]struct{}),
		interval:       24 * time.Hour,
		clearChan:      make(chan bool),
	}
}

func (r *Reminder) Start() error {
	if r.RemindTime.IsZero() {
		return errors.New("no remind time set")
	}
	if r.Started {
		return errors.New("reminder already started")
	}

	r.ticker = time.NewTicker(time.Second)
	r.Started = true

	go func() {
		for {
			select {
			case <-r.ticker.C:
				r.processTick()
			case <-r.clearChan:
				r.processClean()
			}
		}
	}()

	return nil
}

func (r *Reminder) Stop() error {
	r.clearChan <- true
	r.Started = false

	return nil
}

func (r *Reminder) SetRemindTime(t time.Time) {
	if t.Unix() < time.Now().Unix() {
		r.RemindTime = t.Add(r.interval)
		return
	}

	r.RemindTime = t
}

func (r *Reminder) SetWeekdaysToSkip(weekdays []time.Weekday) {
	r.weekdaysToSkip = make(map[time.Weekday]struct{})
	empty := struct{}{}

	for _, weekday := range weekdays {
		r.weekdaysToSkip[weekday] = empty
	}
}

func (r *Reminder) PrintWeekdaysToSkip() string {
	var weekdays []string
	for key := range r.weekdaysToSkip {
		weekdays = append(weekdays, string(key))
	}
	return strings.Join(weekdays, ", ")
}

func (r *Reminder) processTick() {
	now := time.Now()

	if _, skip := r.weekdaysToSkip[now.Weekday()]; skip {
		r.RemindTime = r.RemindTime.Add(r.interval)
		return
	}

	if now.Unix() >= r.RemindTime.Unix() {
		r.RemindChan <- r.RemindMessage
		r.RemindTime = r.RemindTime.Add(r.interval)
	}
}

func (r *Reminder) processClean() {
	r.ticker.Stop()
}
