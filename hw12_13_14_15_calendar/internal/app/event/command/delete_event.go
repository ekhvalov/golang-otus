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
	repository event.Repository
}

func NewDeleteEventRequestHandler(repository event.Repository) (DeleteEventRequestHandler, error) {
	if repository == nil {
		return nil, fmt.Errorf("provided repository is nil")
	}
	return deleteEventRequestHandler{repository: repository}, nil
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
