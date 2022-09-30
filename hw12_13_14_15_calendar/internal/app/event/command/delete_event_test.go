package command_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/event/command"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_deleteEventRequestHandler_Handle_Error(t *testing.T) {
	getMockRepository := func(controller *gomock.Controller) event.Repository {
		return mock.NewMockRepository(controller)
	}
	eventID := "100500"
	errDelete := errors.New("delete event error")
	tests := map[string]struct {
		request       command.DeleteEventRequest
		getRepository func(controller *gomock.Controller) event.Repository
	}{
		"event ID is empty": {
			request:       command.DeleteEventRequest{ID: ""},
			getRepository: getMockRepository,
		},
		"should return error when event repository returned error": {
			request: command.DeleteEventRequest{ID: eventID},
			getRepository: func(controller *gomock.Controller) event.Repository {
				r := mock.NewMockRepository(controller)
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

			h, err := command.NewDeleteEventRequestHandler(tt.getRepository(controller))

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
	repository := mock.NewMockRepository(controller)
	repository.EXPECT().
		Delete(context.Background(), eventID).
		Return(nil)

	h, err := command.NewDeleteEventRequestHandler(repository)

	require.NoError(t, err)
	err = h.Handle(context.Background(), command.DeleteEventRequest{ID: eventID})
	require.NoError(t, err)
}
