package query

import (
	"context"
	"fmt"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
)

type GetDayEventsRequest struct {
	Date time.Time
}

type GetDayEventsResponse struct {
	Events []event.Event
}

type GetDayEventsRequestHandler interface {
	Handle(ctx context.Context, request GetDayEventsRequest) (*GetDayEventsResponse, error)
}

type getDayEventsRequestHandler struct {
	repository event.Repository
}

func (h getDayEventsRequestHandler) Handle(
	ctx context.Context,
	request GetDayEventsRequest,
) (*GetDayEventsResponse, error) {
	events, err := h.repository.GetDayEvents(ctx, request.Date)
	if err != nil {
		return nil, fmt.Errorf("repository GetDayEvents error: %w", err)
	}
	return &GetDayEventsResponse{Events: events}, nil
}
