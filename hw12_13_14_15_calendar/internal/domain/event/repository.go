package event

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, event Event) error
	Update(ctx context.Context, event Event) error
	Delete(ctx context.Context, event Event) error
	GetDayEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetWeekEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetMonthEvents(ctx context.Context, date time.Time) ([]Event, error)
}
