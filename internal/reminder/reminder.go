package reminder

import (
	"time"

	"github.com/pkg/errors"
)

var ticker *time.Ticker
var clearChan chan bool

func NewReminder(remindMessage string) *Reminder {
	return &Reminder{
		RemindMessage: remindMessage,
		Interval:      24 * time.Hour,
	}
}

func (r *Reminder) Start(remindChan chan string) error {
	if r.RemindAt.IsZero() {
		return errors.New("remindAt is not set")
	}
	if r.Started {
		return errors.New("reminder already started")
	}

	ticker = time.NewTicker(time.Second)
	clearChan = make(chan bool)

	r.Started = true

	go func() {
		for {
			select {
			case <-ticker.C:
				if time.Now().Unix() >= r.RemindAt.Unix() {
					r.RemindAt = r.RemindAt.Add(r.Interval)
					remindChan <- r.RemindMessage
				}
			case <-clearChan:
				ticker.Stop()
				return
			}
		}
	}()

	return nil
}

func (r *Reminder) Stop() error {
	clearChan <- true
	r.Started = false

	return nil
}
