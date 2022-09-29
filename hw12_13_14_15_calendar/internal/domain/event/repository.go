package event

//go:generate mockgen -destination=./mock/repository_gen.go -package mock . Repository

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, event Event) error
	Update(ctx context.Context, eventID string, event Event) error
	Delete(ctx context.Context, eventID string) error
	GetDayEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetWeekEvents(ctx context.Context, date time.Time) ([]Event, error)
	GetMonthEvents(ctx context.Context, date time.Time) ([]Event, error)
	IsDateAvailable(ctx context.Context, date time.Time, duration time.Duration) (bool, error)
}
