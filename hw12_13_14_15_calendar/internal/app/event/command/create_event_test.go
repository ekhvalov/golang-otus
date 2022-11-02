package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ekhvalov/otus-golang/hw12_13_14_15_calendar/internal/domain/event"
	storagemock "github.com/ekhvalov/otus-golang/hw12_13_14_15_calendar/internal/domain/event/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	request = CreateEventRequest{
		Title:        "Event 1",
		DateTime:     time.Now().Add(time.Hour),
		Duration:     time.Hour,
		UserID:       "10",
		Description:  "Description",
		NotifyBefore: time.Hour,
	}
	newEvent = event.Event{
		ID:           "10",
		Title:        request.Title,
		DateTime:     request.DateTime,
		Duration:     request.Duration,
		UserID:       request.UserID,
		Description:  request.Description,
		NotifyBefore: request.NotifyBefore,
	}
)

func Test_When_ValidationErrorOccurred_Then_HandlerShouldReturnError(t *testing.T) {
	tests := map[string]struct {
		request CreateEventRequest
	}{
		"empty title": {
			request: CreateEventRequest{
				Title:        "",
				DateTime:     time.Time{},
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
		},
		"time is in the past": {
			request: CreateEventRequest{
				Title:        "Event 1",
				DateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
		},
		"duration is less than minimal duration": {
			request: CreateEventRequest{
				Title:        "Event 1",
				DateTime:     time.Now().Add(time.Hour),
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
		},
		"user ID is empty": {
			request: CreateEventRequest{
				Title:        "Event 1",
				DateTime:     time.Now().Add(time.Hour),
				Duration:     time.Hour,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			c := createEventRequestHandler{storage: nil}
			got, err := c.Handle(context.Background(), tt.request)
			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}

func Test_When_StorageCreateErrorOccurred_Then_HandlerShouldReturnError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	storage := storagemock.NewMockStorage(controller)
	errCreateEvent := errors.New("create event error")
	storage.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(event.Event{}, errCreateEvent)
	h := createEventRequestHandler{storage: storage}

	response, err := h.Handle(context.Background(), request)
	require.Error(t, err)
	require.ErrorIs(t, err, errCreateEvent)
	require.Nil(t, response)
}

func Test_When_NoErrorsOccurred_Then_HandlerShouldReturnCreateEventResponse(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	storage := storagemock.NewMockStorage(controller)

	storage.EXPECT().
		Create(gomock.Any(), event.Event{
			ID:           "",
			Title:        request.Title,
			DateTime:     request.DateTime,
			Duration:     request.Duration,
			UserID:       request.UserID,
			Description:  request.Description,
			NotifyBefore: request.NotifyBefore,
		}).
		Return(newEvent, nil)
	h := createEventRequestHandler{storage: storage}

	response, err := h.Handle(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, &CreateEventResponse{Event: newEvent}, response)
}
