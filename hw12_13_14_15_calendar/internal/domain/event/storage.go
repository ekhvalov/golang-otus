package event

//go:generate mockgen -destination=./mock/storage_gen.go -package mock . Storage

import (
	"context"
	"errors"
	"time"
)

var ErrDateBusy = errors.New("event date is busy")

type Storage interface {
	// Create When event time is overlapped with existed events ErrDateBusy will be returned
	Create(ctx context.Context, event Event) error
	// Update When event time is overlapped with existed events ErrDateBusy will be returned
	Update(ctx context.Context, eventID string, event Event) error
	Delete(ctx context.Context, eventID string) error
	GetDayEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetWeekEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetMonthEvents(ctx context.Context, date time.Time) ([]Event, error)
}
