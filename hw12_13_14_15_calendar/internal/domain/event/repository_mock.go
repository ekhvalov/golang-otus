package event

import (
	"context"
	"fmt"
	"time"
)

type BrokenRepository struct{}

func (r BrokenRepository) Create(_ context.Context, _ Event) error {
	return fmt.Errorf("create error")
}

func (r BrokenRepository) Update(_ context.Context, _ string, _ Event) error {
	return fmt.Errorf("update error")
}

func (r BrokenRepository) Delete(_ context.Context, _ string) error {
	return fmt.Errorf("delete error")
}

func (r BrokenRepository) GetDayEvents(_ context.Context, _ time.Time) ([]Event, error) {
	return nil, fmt.Errorf("get day event error")
}

func (r BrokenRepository) GetWeekEvents(_ context.Context, _ time.Time) ([]Event, error) {
	return nil, fmt.Errorf("get week event error")
}

func (r BrokenRepository) GetMonthEvents(_ context.Context, _ time.Time) ([]Event, error) {
	return nil, fmt.Errorf("get month event error")
}

func (r BrokenRepository) IsDateAvailable(_ context.Context, _ time.Time, _ time.Duration) (bool, error) {
	return false, fmt.Errorf("is timeslot available error")
}

type PlainRepository struct {
	Event       Event
	EventID     string
	Date        time.Time
	Duration    time.Duration
	Events      []Event
	IsAvailable bool
}

func (r *PlainRepository) Create(_ context.Context, event Event) error {
	r.Event = event
	return nil
}

func (r *PlainRepository) Update(_ context.Context, eventID string, event Event) error {
	r.Event = event
	r.EventID = eventID
	return nil
}

func (r *PlainRepository) Delete(_ context.Context, eventID string) error {
	r.EventID = eventID
	return nil
}

func (r *PlainRepository) GetDayEvents(_ context.Context, date time.Time) ([]Event, error) {
	r.Date = date
	return r.Events, nil
}

func (r *PlainRepository) GetWeekEvents(_ context.Context, date time.Time) ([]Event, error) {
	r.Date = date
	return r.Events, nil
}

func (r *PlainRepository) GetMonthEvents(_ context.Context, date time.Time) ([]Event, error) {
	r.Date = date
	return r.Events, nil
}

func (r *PlainRepository) IsDateAvailable(_ context.Context, date time.Time, duration time.Duration) (bool, error) {
	r.Date = date
	r.Duration = duration
	return r.IsAvailable, nil
}
