package query

import (
	"context"
	"testing"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/stretchr/testify/require"
)

func Test_getMonthEventsRequestHandler_Handle(t *testing.T) {
	type fields struct {
		repository event.Repository
	}
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
	tests := map[string]struct {
		fields   fields
		request  GetMonthEventsRequest
		want     *GetMonthEventsResponse
		wantDate time.Time
		wantErr  bool
	}{
		"should return error when repository returned error": {
			fields:  fields{repository: event.BrokenRepository{}},
			request: GetMonthEventsRequest{Date: date},
			wantErr: true,
		},
		"should return events when no error returned by repository": {
			fields:   fields{repository: &event.PlainRepository{Events: events}},
			request:  GetMonthEventsRequest{Date: date},
			wantErr:  false,
			want:     &GetMonthEventsResponse{Events: events},
			wantDate: date,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			h := getMonthEventsRequestHandler{
				repository: tt.fields.repository,
			}
			got, err := h.Handle(context.Background(), tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			r := tt.fields.repository.(*event.PlainRepository)
			require.Equal(t, tt.wantDate, r.Date)
			require.Equal(t, tt.want, got)
		})
	}
}
