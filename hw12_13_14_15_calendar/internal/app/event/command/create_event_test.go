package command

import (
	"context"
	"errors"
	"testing"
	"time"

	providermock "github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/event/command/mock"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	storagemock "github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

var (
	eventID            = "100500"
	getPlainIDProvider = func(controller *gomock.Controller) IDProvider {
		p := providermock.NewMockIDProvider(controller)
		p.EXPECT().GetID().Return(eventID, nil)
		return p
	}
	request = CreateEventRequest{
		Title:        "Event 1",
		DateTime:     time.Now().Add(time.Hour),
		Duration:     time.Hour,
		UserID:       "10",
		Description:  "Description",
		NotifyBefore: time.Hour,
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
			c := createEventRequestHandler{idProvider: nil, storage: nil}
			got, err := c.Handle(context.Background(), tt.request)
			require.Error(t, err)
			require.Nil(t, got)
		})
	}
}

func Test_When_StorageIsAvailableReturnedError_Then_HandlerShouldReturnError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	errStorage := errors.New("storage error")
	storage := storagemock.NewMockStorage(controller)
	storage.EXPECT().
		IsDateAvailable(context.Background(), request.DateTime, request.Duration).
		Return(false, errStorage)
	h := createEventRequestHandler{
		idProvider: nil,
		storage:    storage,
	}

	response, err := h.Handle(context.Background(), request)

	require.Error(t, err)
	require.ErrorIs(t, err, errStorage)
	require.Nil(t, response)
}

func Test_When_RequestedDateIsBusy_Then_HandlerShouldReturnErrDateBusy(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	Storage := storagemock.NewMockStorage(controller)
	Storage.EXPECT().
		IsDateAvailable(context.Background(), request.DateTime, request.Duration).
		Return(false, nil)
	h := createEventRequestHandler{
		idProvider: nil,
		storage:    Storage,
	}

	response, err := h.Handle(context.Background(), request)

	require.Error(t, err)
	require.ErrorIs(t, err, ErrDateBusy)
	require.Nil(t, response)
}

func Test_When_IDProviderReturnedError_Then_HandlerShouldReturnError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	errProvider := errors.New("provider error")
	provider := providermock.NewMockIDProvider(controller)
	provider.EXPECT().GetID().Return("", errProvider)
	Storage := storagemock.NewMockStorage(controller)
	Storage.EXPECT().
		IsDateAvailable(context.Background(), request.DateTime, request.Duration).
		Return(true, nil)
	h := createEventRequestHandler{
		idProvider: provider,
		storage:    Storage,
	}

	response, err := h.Handle(context.Background(), request)
	require.Error(t, err)
	require.ErrorIs(t, err, errProvider)
	require.Nil(t, response)
}

func Test_When_StorageCreateErrorOccurred_Then_HandlerShouldReturnError(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	Storage := storagemock.NewMockStorage(controller)
	Storage.EXPECT().
		IsDateAvailable(context.Background(), request.DateTime, request.Duration).
		Return(true, nil)
	errCreateEvent := errors.New("create event error")
	Storage.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(errCreateEvent)
	h := createEventRequestHandler{
		idProvider: getPlainIDProvider(controller),
		storage:    Storage,
	}

	response, err := h.Handle(context.Background(), request)
	require.Error(t, err)
	require.ErrorIs(t, err, errCreateEvent)
	require.Nil(t, response)
}

func Test_When_NoErrorsOccurred_Then_HandlerShouldReturnCreateEventResponse(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	Storage := storagemock.NewMockStorage(controller)
	Storage.EXPECT().
		IsDateAvailable(context.Background(), request.DateTime, request.Duration).
		Return(true, nil)
	Storage.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(nil)
	h := createEventRequestHandler{
		idProvider: getPlainIDProvider(controller),
		storage:    Storage,
	}

	response, err := h.Handle(context.Background(), request)
	require.NoError(t, err)
	require.Equal(t, &CreateEventResponse{Event: event.Event{
		ID:           eventID,
		Title:        request.Title,
		DateTime:     request.DateTime,
		Duration:     request.Duration,
		UserID:       request.UserID,
		Description:  request.Description,
		NotifyBefore: request.NotifyBefore,
	}}, response)
}
