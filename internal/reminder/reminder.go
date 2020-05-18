package reminder

import (
	"context"
	"go.opencensus.io/trace"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func New(interval time.Duration, location *time.Location) *Reminder {
	return &Reminder{
		WeekdaysToSkip: make(map[time.Weekday]struct{}),
		Interval:       interval,
		Location:       location,
		clearChan:      make(chan struct{}, 1),
	}
}

func (r *Reminder) Start(rawRemindTime string) (chan struct{}, error) {
	_, span := trace.StartSpan(context.Background(), "reminder.Reminder.Start")
	defer span.End()

	remindTime, err := r.parseTime(rawRemindTime)
	if err != nil {
		return nil, errors.Wrap(err, "no remind time set")
	}

	r.RemindTime = remindTime

	if r.Started {
		return r.remindChan, nil
	}

	r.remindChan = make(chan struct{}, 1)
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

	return r.remindChan, nil
}

func (r *Reminder) Stop() error {
	_, span := trace.StartSpan(context.Background(), "reminder.Reminder.Stop")
	defer span.End()

	r.clearChan <- struct{}{}

	return nil
}

func (r *Reminder) SetWeekdaysToSkip(rawWeekdaysToSkip string) error {
	rawDays := strings.Split(rawWeekdaysToSkip, ",")

	for _, rawDay := range rawDays {
		weekday, err := strconv.Atoi(strings.TrimSpace(rawDay))
		if err == nil {
			r.WeekdaysToSkip[time.Weekday(weekday)] = struct{}{}
		}
	}

	return nil
}

func (r *Reminder) processTick() {
	now := time.Now()

	if _, skip := r.WeekdaysToSkip[r.RemindTime.Weekday()]; skip {
		r.RemindTime = r.RemindTime.Add(r.Interval)
		return
	}

	if now.Unix() >= r.RemindTime.Unix() {
		r.remindChan <- struct{}{}
		r.RemindTime = r.RemindTime.Add(r.Interval)
	}
}

func (r *Reminder) processClean() {
	if r.Started {
		defer close(r.remindChan)
		r.ticker.Stop()
		r.Started = false
	}
}

func (r *Reminder) parseTime(rawRemindTime string) (time.Time, error) {
	emptyTime := time.Time{}
	hmArr := strings.Split(rawRemindTime, ":")

	if len(hmArr) != 2 {
		return emptyTime, errors.New("invalid remind time format")
	}

	hour, err := strconv.Atoi(hmArr[0])
	if err != nil {
		return emptyTime, errors.Wrap(err, `invalid remind "hour" value`)
	}

	min, err := strconv.Atoi(hmArr[1])
	if err != nil {
		return emptyTime, errors.Wrap(err, `invalid remind "minute" value`)
	}

	now := time.Now()
	remindTime := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		hour,
		min,
		0,
		0,
		r.Location,
	).UTC()

	for remindTime.Unix() < now.Unix() {
		remindTime = remindTime.Add(r.Interval)
	}

	return remindTime, nil
}
