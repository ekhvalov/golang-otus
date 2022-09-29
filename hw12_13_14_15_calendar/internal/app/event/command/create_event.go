package command

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
)

var ErrDateBusy = errors.New("requested date is busy")

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

type createEventRequestHandler struct {
	idProvider IDProvider
	repository event.Repository
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
	isAvailable, err := h.repository.IsDateAvailable(ctx, request.DateTime, request.Duration)
	if err != nil {
		return nil, err
	}
	if !isAvailable {
		return nil, ErrDateBusy
	}
	ID, err := h.idProvider.GetID()
	if err != nil {
		return nil, fmt.Errorf("provide ID error: %w", err)
	}
	e := event.Event{
		ID:           ID,
		Title:        request.Title,
		DateTime:     request.DateTime,
		Duration:     request.Duration,
		UserID:       request.UserID,
		Description:  request.Description,
		NotifyBefore: request.NotifyBefore,
	}
	err = h.repository.Create(ctx, e)
	if err != nil {
		return nil, fmt.Errorf("repository create event error: %w", err)
	}
	return &CreateEventResponse{Event: e}, nil
}
