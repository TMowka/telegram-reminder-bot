package reminder

import (
	"context"
	"go.opencensus.io/trace"
	"time"

	"github.com/pkg/errors"
)

func New(interval time.Duration) *Reminder {
	return &Reminder{
		WeekdaysToSkip: make(map[time.Weekday]struct{}),
		Interval:       interval,
		clearChan:      make(chan struct{}, 1),
	}
}

func (r *Reminder) Start(remindTime time.Time) (chan struct{}, error) {
	_, span := trace.StartSpan(context.Background(), "reminder.Reminder.Start")
	defer span.End()

	if remindTime.IsZero() {
		return nil, errors.New("no remind time set")
	}

	for remindTime.Unix() < time.Now().Unix() {
		remindTime = remindTime.Add(r.Interval)
	}

	r.remindTime = remindTime

	if r.started {
		return r.remindChan, nil
	}

	r.remindChan = make(chan struct{}, 1)
	r.ticker = time.NewTicker(time.Second)
	r.started = true

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

	return r.remindChan, nil
}

func (r *Reminder) Stop() error {
	_, span := trace.StartSpan(context.Background(), "reminder.Reminder.Stop")
	defer span.End()

	r.clearChan <- struct{}{}

	return nil
}

func (r *Reminder) processTick() {
	now := time.Now()

	if _, skip := r.WeekdaysToSkip[r.remindTime.Weekday()]; skip {
		r.remindTime = r.remindTime.Add(r.Interval)
		return
	}

	if now.Unix() >= r.remindTime.Unix() {
		r.remindChan <- struct{}{}
		r.remindTime = r.remindTime.Add(r.Interval)
	}
}

func (r *Reminder) processClean() {
	defer close(r.remindChan)

	r.ticker.Stop()
	r.started = false
}
