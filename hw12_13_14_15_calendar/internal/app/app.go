package app

import (
	"context"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
)

type App struct {
	logger  Logger
	storage event.Storage
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func New(logger Logger, storage event.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.Create(ctx, event.Event{ID: id, Title: title})
}
