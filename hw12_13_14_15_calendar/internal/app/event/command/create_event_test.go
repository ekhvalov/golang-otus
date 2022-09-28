package command

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/stretchr/testify/require"
)

type plainIDProvider struct {
	id string
}

func (i plainIDProvider) GetID() (string, error) {
	return i.id, nil
}

type brokenIDProvider struct{}

func (b brokenIDProvider) GetID() (string, error) {
	return "", fmt.Errorf("ID provide error")
}

func Test_createEventRequestHandler_Handle(t *testing.T) {
	type fields struct {
		idProvider IDProvider
		repository event.Repository
	}
	eventTime := time.Now().Add(time.Hour)
	eventDuration := time.Hour

	tests := map[string]struct {
		fields  fields
		request CreateEventRequest
		want    *CreateEventResponse
		wantErr bool
	}{
		"should return error when empty title provided": {
			request: CreateEventRequest{
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
			request: CreateEventRequest{
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
			request: CreateEventRequest{
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
			request: CreateEventRequest{
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     time.Hour,
				UserID:       "",
				Description:  "",
				NotifyBefore: 0,
			},
			wantErr: true,
		},
		"should return error when IDProvider returned error": {
			request: CreateEventRequest{
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     eventDuration,
				UserID:       "10",
				Description:  "",
				NotifyBefore: 0,
			},
			fields: fields{
				idProvider: brokenIDProvider{},
				repository: nil,
			},
			wantErr: true,
		},
		"should return error when event.Repository returned error": {
			request: CreateEventRequest{
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     eventDuration,
				UserID:       "10",
				Description:  "",
				NotifyBefore: 0,
			},
			fields: fields{
				idProvider: plainIDProvider{id: "100500"},
				repository: event.BrokenRepository{},
			},
			wantErr: true,
		},
		"should return event.Event when no errors found": {
			request: CreateEventRequest{
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     eventDuration,
				UserID:       "10",
				Description:  "Description",
				NotifyBefore: time.Hour,
			},
			fields: fields{
				idProvider: plainIDProvider{id: "100500"},
				repository: &event.PlainRepository{},
			},
			wantErr: false,
			want: &CreateEventResponse{Event: event.Event{
				ID:           "100500",
				Title:        "Event 1",
				DateTime:     eventTime,
				Duration:     eventDuration,
				UserID:       "10",
				Description:  "Description",
				NotifyBefore: time.Hour,
			}},
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			c := createEventRequestHandler{
				idProvider: tt.fields.idProvider,
				repository: tt.fields.repository,
			}
			got, err := c.Handle(context.Background(), tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			r := tt.fields.repository.(*event.PlainRepository)
			require.Equal(t, tt.want.Event, r.Event)
			require.Equal(t, tt.want, got)
		})
	}
}
