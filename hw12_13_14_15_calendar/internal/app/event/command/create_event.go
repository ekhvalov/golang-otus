package command

import (
	"context"
	"fmt"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
)

type CreateEventRequest struct {
	Title        string
	DateTime     time.Time
	Duration     time.Duration
	UserID       string
	Description  string
	NotifyBefore time.Duration
}

type CreateEventResponse struct {
	Event event.Event
}

type CreateEventRequestHandler interface {
	Handle(ctx context.Context, request CreateEventRequest) (*CreateEventResponse, error)
}

func NewCreateEventRequestHandler(storage event.Storage) (CreateEventRequestHandler, error) {
	if storage == nil {
		return nil, fmt.Errorf("storage is nil")
	}
	return &createEventRequestHandler{storage: storage}, nil
}

type createEventRequestHandler struct {
	storage event.Storage
}

func (h createEventRequestHandler) Handle(
	ctx context.Context,
	request CreateEventRequest,
) (*CreateEventResponse, error) {
	if err := validateTitle(request.Title); err != nil {
		return nil, err
	}
	if err := validateDateTime(request.DateTime); err != nil {
		return nil, err
	}
	if err := validateDuration(request.Duration); err != nil {
		return nil, err
	}
	if err := validateUserID(request.UserID); err != nil {
		return nil, err
	}
	e := event.Event{
		Title:        request.Title,
		DateTime:     request.DateTime,
		Duration:     request.Duration,
		UserID:       request.UserID,
		Description:  request.Description,
		NotifyBefore: request.NotifyBefore,
	}
	newEvent, err := h.storage.Create(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("create event error: %w", err)
	}
	return &CreateEventResponse{Event: newEvent}, nil
}
