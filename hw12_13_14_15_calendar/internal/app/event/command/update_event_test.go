package command

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_updateEventRequestHandler_Handle(t *testing.T) {
	request := UpdateEventRequest{
		ID:           "100500",
		Title:        "Title",
		DateTime:     time.Now().Add(time.Hour),
		Duration:     time.Hour,
		UserID:       "100600",
		Description:  "Description",
		NotifyBefore: time.Hour,
	}
	controller := gomock.NewController(t)
	defer controller.Finish()

	storage := mock.NewMockStorage(controller)
	storage.EXPECT().
		Update(context.Background(), request.ID, event.Event{
			ID:           "",
			Title:        request.Title,
			DateTime:     request.DateTime,
			Duration:     request.Duration,
			UserID:       request.UserID,
			Description:  request.Description,
			NotifyBefore: request.NotifyBefore,
		}).
		Return(nil)

	h := updateEventRequestHandler{storage: storage}

	err := h.Handle(context.Background(), request)
	require.NoError(t, err)
}

func Test_updateEventRequestHandler_Handle_Error(t *testing.T) {
	request := UpdateEventRequest{
		ID:           "100500",
		Title:        "Title",
		DateTime:     time.Now().Add(time.Hour),
		Duration:     time.Hour,
		UserID:       "100600",
		Description:  "Description",
		NotifyBefore: time.Hour,
	}
	errUpdate := errors.New("update event error")
	getMockStorage := func(controller *gomock.Controller) event.Storage {
		return mock.NewMockStorage(controller)
	}
	tests := map[string]struct {
		request    UpdateEventRequest
		getStorage func(controller *gomock.Controller) event.Storage
		wantErr    error
	}{
		"validation error when empty ID provided": {
			request: UpdateEventRequest{
				ID:           "",
				Title:        "",
				DateTime:     time.Time{},
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			getStorage: getMockStorage,
		},
		"validation error when empty title provided": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "",
				DateTime:     time.Time{},
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			getStorage: getMockStorage,
		},
		"validation error when time is in the past": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "Event 1",
				DateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			getStorage: getMockStorage,
		},
		"validation error when duration is less than minimal duration": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "Event 1",
				DateTime:     request.DateTime,
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			getStorage: getMockStorage,
		},
		"validation error when user ID is empty": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "Event 1",
				DateTime:     request.DateTime,
				Duration:     time.Hour,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			getStorage: getMockStorage,
		},
		"storage.Update error": {
			request: request,
			getStorage: func(controller *gomock.Controller) event.Storage {
				r := mock.NewMockStorage(controller)
				r.EXPECT().
					Update(context.Background(), request.ID, event.Event{
						Title:        request.Title,
						DateTime:     request.DateTime,
						Duration:     request.Duration,
						UserID:       request.UserID,
						Description:  request.Description,
						NotifyBefore: request.NotifyBefore,
					}).
					Return(errUpdate)
				return r
			},
			wantErr: errUpdate,
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			h := updateEventRequestHandler{storage: tt.getStorage(controller)}

			err := h.Handle(context.Background(), tt.request)

			require.Error(t, err)
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			}
		})
	}
}
