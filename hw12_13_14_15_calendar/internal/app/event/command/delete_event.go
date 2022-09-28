package command

import (
	"context"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
)

type DeleteEventRequest struct {
	ID string
}

type DeleteEventRequestHandler interface {
	Handle(ctx context.Context, request DeleteEventRequest) error
}

type deleteEventRequestHandler struct {
	repository event.Repository
}

func (h deleteEventRequestHandler) Handle(ctx context.Context, request DeleteEventRequest) error {
	if err := validateID(request.ID); err != nil {
		return err
	}
	err := h.repository.Delete(ctx, request.ID)
	if err != nil {
		return err
	}
	return nil
}
