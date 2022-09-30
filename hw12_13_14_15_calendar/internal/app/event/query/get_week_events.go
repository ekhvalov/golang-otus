package query

import (
	"context"
	"fmt"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
)

type GetWeekEventsRequest struct {
	Date time.Time
}

type GetWeekEventsResponse struct {
	Events []event.Event
}

type GetWeekEventsRequestHandler interface {
	Handle(ctx context.Context, request GetWeekEventsRequest) (*GetWeekEventsResponse, error)
}

func NewGetWeekEventsRequestHandler(repository event.Repository) (GetWeekEventsRequestHandler, error) {
	if repository == nil {
		return nil, fmt.Errorf("provided repository is nil")
	}
	return getWeekEventsRequestHandler{repository: repository}, nil
}

type getWeekEventsRequestHandler struct {
	repository event.Repository
}

func (h getWeekEventsRequestHandler) Handle(
	ctx context.Context,
	request GetWeekEventsRequest,
) (*GetWeekEventsResponse, error) {
	events, err := h.repository.GetWeekEvents(ctx, request.Date)
	if err != nil {
		return nil, fmt.Errorf("repository GetWeekEvents error: %w", err)
	}
	return &GetWeekEventsResponse{Events: events}, nil
}
