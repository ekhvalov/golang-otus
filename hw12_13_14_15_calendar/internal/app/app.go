package app

import (
	"context"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, e storage.Event) error
	UpdateEvent(ctx context.Context, e storage.Event) error
	DeleteEvent(ctx context.Context, e storage.Event) error
	GetEvents(ctx context.Context) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(ctx, storage.Event{ID: id, Title: title})
}
