package command

import (
	"context"
	"fmt"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
)

type DeleteEventRequest struct {
	ID string
}

type DeleteEventRequestHandler interface {
	Handle(ctx context.Context, request DeleteEventRequest) error
}

type deleteEventRequestHandler struct {
	storage event.Storage
}

func NewDeleteEventRequestHandler(storage event.Storage) (DeleteEventRequestHandler, error) {
	if storage == nil {
		return nil, fmt.Errorf("provided storage is nil")
	}
	return deleteEventRequestHandler{storage: storage}, nil
}

func (h deleteEventRequestHandler) Handle(ctx context.Context, request DeleteEventRequest) error {
	if err := validateID(request.ID); err != nil {
		return err
	}
	err := h.storage.Delete(ctx, request.ID)
	if err != nil {
		return err
	}
	return nil
}
