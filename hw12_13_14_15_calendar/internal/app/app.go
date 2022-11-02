package app

//go:generate mockgen -destination=./mock/app.gen.go -package mock . Application
//go:generate mockgen -destination=./mock/logger.gen.go -package mock . Logger

import (
	"context"
	"errors"
	"fmt"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/command"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/query"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
)

type Application interface {
	CreateEvent(request command.CreateEventRequest) (*command.CreateEventResponse, error)
	UpdateEvent(request command.UpdateEventRequest) error
	DeleteEvent(request command.DeleteEventRequest) error
	GetDayEvents(request query.GetDayEventsRequest) (*query.GetDayEventsResponse, error)
	GetWeekEvents(request query.GetWeekEventsRequest) (*query.GetWeekEventsResponse, error)
	GetMonthEvents(request query.GetMonthEventsRequest) (*query.GetMonthEventsResponse, error)
}

type app struct {
	createHandler   command.CreateEventRequestHandler
	updateHandler   command.UpdateEventRequestHandler
	deleteHandler   command.DeleteEventRequestHandler
	getDayHandler   query.GetDayEventsRequestHandler
	getWeekHandler  query.GetWeekEventsRequestHandler
	getMonthHandler query.GetMonthEventsRequestHandler
	logger          Logger
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func New(logger Logger, storage event.Storage) (Application, error) {
	createHandler, err := command.NewCreateEventRequestHandler(storage)
	if err != nil {
		return nil, fmt.Errorf("create CreateEventRequestHandler error: %w", err)
	}
	updateHandler, err := command.NewUpdateEventRequestHandler(storage)
	if err != nil {
		return nil, fmt.Errorf("create UpdateEventRequestHandler error: %w", err)
	}
	deleteHandler, err := command.NewDeleteEventRequestHandler(storage)
	if err != nil {
		return nil, fmt.Errorf("create DeleteEventRequestHandler error: %w", err)
	}
	getDayHandler, err := query.NewGetDayEventsRequestHandler(storage)
	if err != nil {
		return nil, fmt.Errorf("create GetDayEventsRequestHandler error: %w", err)
	}
	getWeekHandler, err := query.NewGetWeekEventsRequestHandler(storage)
	if err != nil {
		return nil, fmt.Errorf("create GetWeekEventsRequestHandler error: %w", err)
	}
	getMonthHandler, err := query.NewGetMonthEventsRequestHandler(storage)
	if err != nil {
		return nil, fmt.Errorf("create GetMonthEventsRequestHandler error: %w", err)
	}
	return &app{
		createHandler:   createHandler,
		updateHandler:   updateHandler,
		deleteHandler:   deleteHandler,
		getDayHandler:   getDayHandler,
		getWeekHandler:  getWeekHandler,
		getMonthHandler: getMonthHandler,
		logger:          logger,
	}, nil
}

func (a *app) CreateEvent(request command.CreateEventRequest) (*command.CreateEventResponse, error) {
	response, err := a.createHandler.Handle(context.Background(), request)
	if err != nil {
		if errors.Is(err, event.ErrStorage{}) {
			a.logger.Error(fmt.Sprintf("create event handler storage error: %s", err))
		}
		return nil, err
	}
	return response, nil
}

func (a *app) UpdateEvent(request command.UpdateEventRequest) error {
	err := a.updateHandler.Handle(context.Background(), request)
	if err != nil {
		if errors.Is(err, event.ErrStorage{}) {
			a.logger.Error(fmt.Sprintf("update event handler storage error: %s", err))
		}
		return err
	}
	return nil
}

func (a *app) DeleteEvent(request command.DeleteEventRequest) error {
	err := a.deleteHandler.Handle(context.Background(), request)
	if err != nil {
		if errors.Is(err, event.ErrStorage{}) {
			a.logger.Error(fmt.Sprintf("delete event handler storage error: %s", err))
		}
		return err
	}
	return nil
}

func (a *app) GetDayEvents(request query.GetDayEventsRequest) (*query.GetDayEventsResponse, error) {
	response, err := a.getDayHandler.Handle(context.Background(), request)
	if err != nil {
		if errors.Is(err, event.ErrStorage{}) {
			a.logger.Error(fmt.Sprintf("get day events handler storage error: %s", err))
		}
		return nil, err
	}
	return response, nil
}

func (a *app) GetWeekEvents(request query.GetWeekEventsRequest) (*query.GetWeekEventsResponse, error) {
	response, err := a.getWeekHandler.Handle(context.Background(), request)
	if err != nil {
		if errors.Is(err, event.ErrStorage{}) {
			a.logger.Error(fmt.Sprintf("get week events handler storage error: %s", err))
		}
		return nil, err
	}
	return response, nil
}

func (a *app) GetMonthEvents(request query.GetMonthEventsRequest) (*query.GetMonthEventsResponse, error) {
	response, err := a.getMonthHandler.Handle(context.Background(), request)
	if err != nil {
		if errors.Is(err, event.ErrStorage{}) {
			a.logger.Error(fmt.Sprintf("get month events handler storage error: %s", err))
		}
		return nil, err
	}
	return response, nil
}
