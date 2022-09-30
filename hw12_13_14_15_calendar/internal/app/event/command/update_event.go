package command

import (
	"context"
	"fmt"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
)

type UpdateEventRequest struct {
	ID           string
	Title        string
	DateTime     time.Time
	Duration     time.Duration
	UserID       string
	Description  string
	NotifyBefore time.Duration
}

type UpdateEventRequestHandler interface {
	Handle(ctx context.Context, request UpdateEventRequest) error
}

type updateEventRequestHandler struct {
	repository event.Repository
}

func (h updateEventRequestHandler) Handle(ctx context.Context, request UpdateEventRequest) error {
	if err := validateID(request.ID); err != nil {
		return err
	}
	if err := validateTitle(request.Title); err != nil {
		return err
	}
	if err := validateDateTime(request.DateTime); err != nil {
		return err
	}
	if err := validateDuration(request.Duration); err != nil {
		return err
	}
	if err := validateUserID(request.UserID); err != nil {
		return err
	}
	isDateAvailable, err := h.repository.IsDateAvailable(ctx, request.DateTime, request.Duration)
	if err != nil {
		return err
	}
	if !isDateAvailable {
		return ErrDateBusy
	}
	err = h.repository.Update(ctx, request.ID, event.Event{
		Title:        request.Title,
		DateTime:     request.DateTime,
		Duration:     request.Duration,
		UserID:       request.UserID,
		Description:  request.Description,
		NotifyBefore: request.NotifyBefore,
	})
	if err != nil {
		return fmt.Errorf("repository update event error: %w", err)
	}
	return nil
}
