package query

import (
	"context"
	"fmt"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
)

type GetMonthEventsRequest struct {
	Date time.Time
}

type GetMonthEventsResponse struct {
	Events []event.Event
}

type GetMonthEventsRequestHandler interface {
	Handle(ctx context.Context, request GetMonthEventsRequest) (*GetMonthEventsResponse, error)
}

type getMonthEventsRequestHandler struct {
	repository event.Repository
}

func (h getMonthEventsRequestHandler) Handle(
	ctx context.Context,
	request GetMonthEventsRequest,
) (*GetMonthEventsResponse, error) {
	events, err := h.repository.GetMonthEvents(ctx, request.Date)
	if err != nil {
		return nil, fmt.Errorf("repository GetMonthEvents error: %w", err)
	}
	return &GetMonthEventsResponse{Events: events}, nil
}
