package event

//go:generate mockgen -destination=./mock/storage_gen.go -package mock . Storage

import (
	"context"
	"errors"
	"time"
)

var (
	ErrDateBusy      = errors.New("event date is busy")
	ErrEventNotFound = errors.New("event not found")
)

// ErrStorage will be returned in case of internal storage error.
type ErrStorage struct {
	message string
}

func (e ErrStorage) Error() string {
	return e.message
}

// NewErrStorage function creates ErrStorage error.
func NewErrStorage(message string) ErrStorage {
	return ErrStorage{message: message}
}

type Storage interface {
	// A Create method stores new Event into a storage
	// Event.ID will be generated inside a storage.
	// When Event.DateTime and Event.Duration are overlapped with existed events ErrDateBusy will be returned.
	// Returns Event with filled ID or ErrStorage if internal storage error occurred.
	Create(ctx context.Context, event Event) (Event, error)
	// An Update method updates an Event.
	// When Event.DateTime and Event.Duration are overlapped with existed events ErrDateBusy will be returned.
	// ErrStorage could be returned if internal storage error occurred.
	Update(ctx context.Context, eventID string, event Event) error
	// A Delete method removes an Event from a storage.
	// ErrStorage could be returned if internal storage error occurred.
	Delete(ctx context.Context, eventID string) error
	// A GetDayEvents method returns a list of Events for a day-long period that starts from the date value.
	GetDayEvents(ctx context.Context, date time.Time) ([]Event, error)
	// A GetWeekEvents method returns a list of Events for a week-long period that starts from the date value.
	GetWeekEvents(ctx context.Context, date time.Time) ([]Event, error)
	// A GetMonthEvents method returns a list of Events for a month-long period that starts from the date value.
	GetMonthEvents(ctx context.Context, date time.Time) ([]Event, error)
}
