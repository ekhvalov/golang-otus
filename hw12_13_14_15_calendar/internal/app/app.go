package app

import (
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
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
