package query

import (
	"context"
	"fmt"
	"time"

	"github.com/ekhvalov/otus-golang/hw12_13_14_15_calendar/internal/domain/event"
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

func NewGetDayEventsRequestHandler(storage event.Storage) (GetDayEventsRequestHandler, error) {
	if storage == nil {
		return nil, fmt.Errorf("provided storage is nil")
	}
	return getDayEventsRequestHandler{storage: storage}, nil
}

type getDayEventsRequestHandler struct {
	storage event.Storage
}

func (h getDayEventsRequestHandler) Handle(
	ctx context.Context,
	request GetDayEventsRequest,
) (*GetDayEventsResponse, error) {
	events, err := h.storage.GetDayEvents(ctx, request.Date)
	if err != nil {
		return nil, fmt.Errorf("storage GetDayEvents error: %w", err)
	}
	return &GetDayEventsResponse{Events: events}, nil
}
