package query

import (
	"context"
	"fmt"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
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

func NewGetWeekEventsRequestHandler(storage event.Storage) (GetWeekEventsRequestHandler, error) {
	if storage == nil {
		return nil, fmt.Errorf("provided storage is nil")
	}
	return getWeekEventsRequestHandler{storage: storage}, nil
}

type getWeekEventsRequestHandler struct {
	storage event.Storage
}

func (h getWeekEventsRequestHandler) Handle(
	ctx context.Context,
	request GetWeekEventsRequest,
) (*GetWeekEventsResponse, error) {
	events, err := h.storage.GetWeekEvents(ctx, request.Date)
	if err != nil {
		return nil, fmt.Errorf("storage GetWeekEvents error: %w", err)
	}
	return &GetWeekEventsResponse{Events: events}, nil
}
