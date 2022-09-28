package command

import (
	"context"
	"testing"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/stretchr/testify/require"
)

func Test_updateEventRequestHandler_Handle(t *testing.T) {
	type fields struct {
		repository event.Repository
	}
	eventTime := time.Now().Add(time.Hour)
	eventDuration := time.Hour

	tests := map[string]struct {
		fields      fields
		request     UpdateEventRequest
		wantErr     bool
		wantEvent   event.Event
		wantEventID string
	}{
		"should return error when empty ID provided": {
			request: UpdateEventRequest{
				ID:           "",
				Title:        "",
				DateTime:     time.Time{},
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			wantErr: true,
		},
		"should return error when empty title provided": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "",
				DateTime:     time.Time{},
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			wantErr: true,
		},
		"should return error when time is in the past": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "Event 1",
				DateTime:     time.Date(2020, 1, 1, 0, 0, 0, 0, time.Local),
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			wantErr: true,
		},
		"should return error when duration is less than minimal duration": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     0,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			wantErr: true,
		},
		"should return error when user ID is empty": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     time.Hour,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			wantErr: true,
		},
		"should return error when event.Repository returned error": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     eventDuration,
				UserID:       "10",
				Description:  "",
				NotifyBefore: 0,
			},
			fields: fields{
				repository: event.BrokenRepository{},
			},
			wantErr: true,
		},
		"should return no error when event fields are correct": {
			request: UpdateEventRequest{
				ID:           "100500",
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     eventDuration,
				UserID:       "100500",
				Description:  "Description",
				NotifyBefore: time.Hour,
			},
			fields: fields{
				repository: &event.PlainRepository{},
			},
			wantErr: false,
			wantEvent: event.Event{
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     eventDuration,
				UserID:       "100500",
				Description:  "Description",
				NotifyBefore: time.Hour,
			},
			wantEventID: "100500",
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			u := updateEventRequestHandler{repository: tt.fields.repository}
			err := u.Handle(context.Background(), tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			r := tt.fields.repository.(*event.PlainRepository)
			require.Equal(t, tt.wantEventID, r.EventID)
			require.Equal(t, tt.wantEvent, r.Event)
		})
	}
}
