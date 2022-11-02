package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ekhvalov/otus-golang/hw12_13_14_15_calendar/internal/app/event/command"
	"github.com/ekhvalov/otus-golang/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/otus-golang/hw12_13_14_15_calendar/internal/domain/event/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_deleteEventRequestHandler_Handle_Error(t *testing.T) {
	getMockStorage := func(controller *gomock.Controller) event.Storage {
		return mock.NewMockStorage(controller)
	}
	eventID := "100500"
	errDelete := errors.New("delete event error")
	tests := map[string]struct {
		request    command.DeleteEventRequest
		getStorage func(controller *gomock.Controller) event.Storage
	}{
		"event ID is empty": {
			request:    command.DeleteEventRequest{ID: ""},
			getStorage: getMockStorage,
		},
		"should return error when event Storage returned error": {
			request: command.DeleteEventRequest{ID: eventID},
			getStorage: func(controller *gomock.Controller) event.Storage {
				r := mock.NewMockStorage(controller)
				r.EXPECT().
					Delete(context.Background(), eventID).
					Return(errDelete)
				return r
			},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()

			h, err := command.NewDeleteEventRequestHandler(tt.getStorage(controller))

			require.NoError(t, err)
			err = h.Handle(context.Background(), tt.request)
			require.Error(t, err)
		})
	}
}

func Test_deleteEventRequestHandler_Handle(t *testing.T) {
	eventID := "100500"
	controller := gomock.NewController(t)
	defer controller.Finish()
	Storage := mock.NewMockStorage(controller)
	Storage.EXPECT().
		Delete(context.Background(), eventID).
		Return(nil)

	h, err := command.NewDeleteEventRequestHandler(Storage)

	require.NoError(t, err)
	err = h.Handle(context.Background(), command.DeleteEventRequest{ID: eventID})
	require.NoError(t, err)
}
