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

func NewGetMonthEventsRequestHandler(storage event.Storage) (GetMonthEventsRequestHandler, error) {
	if storage == nil {
		return nil, fmt.Errorf("provided storage is nil")
	}
	return getMonthEventsRequestHandler{storage: storage}, nil
}

type getMonthEventsRequestHandler struct {
	storage event.Storage
}

func (h getMonthEventsRequestHandler) Handle(
	ctx context.Context,
	request GetMonthEventsRequest,
) (*GetMonthEventsResponse, error) {
	events, err := h.storage.GetMonthEvents(ctx, request.Date)
	if err != nil {
		return nil, fmt.Errorf("storage GetMonthEvents error: %w", err)
	}
	return &GetMonthEventsResponse{Events: events}, nil
}
