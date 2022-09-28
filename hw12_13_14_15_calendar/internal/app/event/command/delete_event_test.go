package command

import (
	"context"
	"testing"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/stretchr/testify/require"
)

func Test_deleteEventRequestHandler_Handle(t *testing.T) {
	type fields struct {
		repository event.Repository
	}
	tests := map[string]struct {
		fields      fields
		request     DeleteEventRequest
		wantErr     bool
		wantEventID string
	}{
		"should return error when event ID is empty": {
			request: DeleteEventRequest{ID: ""},
			wantErr: true,
		},
		"should return error when event repository returned error": {
			request: DeleteEventRequest{ID: "100500"},
			fields:  fields{repository: event.BrokenRepository{}},
			wantErr: true,
		},
		"should return no error when everything is correct": {
			request:     DeleteEventRequest{ID: "100500"},
			fields:      fields{repository: &event.PlainRepository{}},
			wantErr:     false,
			wantEventID: "100500",
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			d := deleteEventRequestHandler{
				repository: tt.fields.repository,
			}
			err := d.Handle(context.Background(), tt.request)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			r := tt.fields.repository.(*event.PlainRepository)
			require.Equal(t, tt.wantEventID, r.EventID)
		})
	}
}
