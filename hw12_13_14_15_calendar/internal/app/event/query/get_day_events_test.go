package query_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/event/query"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_getDayEventsRequestHandler_Handle(t *testing.T) {
	date := time.Now()
	events := []event.Event{
		{
			ID:           "100500",
			Title:        "Title 1",
			DateTime:     date,
			Duration:     time.Hour,
			UserID:       "100500",
			Description:  "Description 1",
			NotifyBefore: time.Hour,
		},
		{
			ID:           "100600",
			Title:        "Title 2",
			DateTime:     date,
			Duration:     time.Hour,
			UserID:       "100600",
			Description:  "Description 2",
			NotifyBefore: time.Hour,
		},
	}
	errGetEvents := errors.New("get events error")

	tests := map[string]struct {
		getRepository func(controller *gomock.Controller) event.Repository
		request       query.GetDayEventsRequest
		want          *query.GetDayEventsResponse
		wantErr       bool
		wantErrType   error
	}{
		"should return error when repository returned error": {
			getRepository: func(controller *gomock.Controller) event.Repository {
				r := mock.NewMockRepository(controller)
				r.EXPECT().
					GetDayEvents(context.Background(), date).
					Return(nil, errGetEvents)
				return r
			},
			request:     query.GetDayEventsRequest{Date: date},
			wantErr:     true,
			wantErrType: errGetEvents,
		},
		"should return events when no error returned by repository": {
			getRepository: func(controller *gomock.Controller) event.Repository {
				r := mock.NewMockRepository(controller)
				r.EXPECT().
					GetDayEvents(context.Background(), date).
					Return(events, nil)
				return r
			},
			request: query.GetDayEventsRequest{Date: date},
			wantErr: false,
			want:    &query.GetDayEventsResponse{Events: events},
		},
	}

	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			h, err := query.NewGetDayEventsRequestHandler(tt.getRepository(controller))
			require.NoError(t, err)

			got, err := h.Handle(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrType != nil {
					require.ErrorIs(t, err, tt.wantErrType)
				}
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
