package internalgrpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/command"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/query"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/mock"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	internalgrpc "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/server/grpc"
	grpcapi "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/pkg/api/grpc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_server_CreateEvent(t *testing.T) {
	eventTime := time.Now()
	tests := map[string]struct {
		getApp  func(controller *gomock.Controller) app.Application
		request *grpcapi.CreateEventRequest
		want    *grpcapi.CreateEventResponse
		wantErr bool
	}{
		"when error occurred then should return error": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					CreateEvent(gomock.Any(), gomock.Any()).
					Return(nil, command.ErrValidate{})
				return a
			},
			wantErr: true,
		},
		"when no error occurred then should return CreateEventResponse": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					CreateEvent(gomock.Any(), gomock.Any()).
					Return(&command.CreateEventResponse{Event: event.Event{
						ID:           "10",
						Title:        "Title",
						DateTime:     eventTime,
						Duration:     time.Minute * 30,
						UserID:       "10",
						Description:  "Description",
						NotifyBefore: time.Minute * 20,
					}}, nil)
				return a
			},
			wantErr: false,
			want:    &grpcapi.CreateEventResponse{Id: "10"},
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			logger := mock.NewMockLogger(controller)
			s, err := internalgrpc.NewServer(tt.getApp(controller), logger)
			require.NoError(t, err)
			srv := s.(grpcapi.CalendarServer)

			createEventResponse, err := srv.CreateEvent(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, createEventResponse)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, createEventResponse)
		})
	}
}

func Test_server_UpdateEvent(t *testing.T) {
	tests := map[string]struct {
		getApp  func(controller *gomock.Controller) app.Application
		request *grpcapi.UpdateEventRequest
		wantErr bool
	}{
		"when error occurred then should return error": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					UpdateEvent(gomock.Any(), gomock.Any()).
					Return(command.ErrValidate{})
				return a
			},
			wantErr: true,
		},
		"when no error occurred then should return Empty": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					UpdateEvent(gomock.Any(), gomock.Any()).
					Return(nil)
				return a
			},
			wantErr: false,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			logger := mock.NewMockLogger(controller)
			s, err := internalgrpc.NewServer(tt.getApp(controller), logger)
			require.NoError(t, err)
			srv := s.(grpcapi.CalendarServer)

			updateEventResponse, err := srv.UpdateEvent(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, updateEventResponse)
				return
			}
			require.NoError(t, err)
			require.Equal(t, &grpcapi.Empty{}, updateEventResponse)
		})
	}
}

func Test_server_DeleteEvent(t *testing.T) {
	tests := map[string]struct {
		getApp  func(controller *gomock.Controller) app.Application
		request *grpcapi.DeleteEventRequest
		wantErr bool
	}{
		"when error occurred then should return error": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					DeleteEvent(gomock.Any(), gomock.Any()).
					Return(event.ErrEventNotFound)
				return a
			},
			wantErr: true,
		},
		"when no error occurred then should return Empty": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					DeleteEvent(gomock.Any(), gomock.Any()).
					Return(nil)
				return a
			},
			wantErr: false,
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			logger := mock.NewMockLogger(controller)
			s, err := internalgrpc.NewServer(tt.getApp(controller), logger)
			require.NoError(t, err)
			srv := s.(grpcapi.CalendarServer)

			deleteEventResponse, err := srv.DeleteEvent(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, deleteEventResponse)
				return
			}
			require.NoError(t, err)
			require.Equal(t, &grpcapi.Empty{}, deleteEventResponse)
		})
	}
}

func Test_server_GetEvents(t *testing.T) {
	tests := map[string]struct {
		getApp  func(controller *gomock.Controller) app.Application
		request *grpcapi.GetEventsRequest
		wantErr bool
		want    *grpcapi.GetEventsResponse
	}{
		"when error occurred then should return error": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetDayEvents(gomock.Any(), gomock.Any()).
					Return(nil, event.ErrStorage{})
				return a
			},
			request: &grpcapi.GetEventsRequest{
				Period: grpcapi.GetEventsRequest_GET_EVENTS_PERIOD_DAY,
			},
			wantErr: true,
		},
		"when events period is not specified then should return error": {
			getApp: func(controller *gomock.Controller) app.Application {
				return mock.NewMockApplication(controller)
			},
			request: &grpcapi.GetEventsRequest{
				Period: grpcapi.GetEventsRequest_GET_EVENTS_PERIOD_UNSPECIFIED,
			},
			wantErr: true,
			want:    &grpcapi.GetEventsResponse{Events: make([]*grpcapi.Event, 0)},
		},
		"when no error occurred then should return GetEventsResponse": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetDayEvents(gomock.Any(), gomock.Any()).
					Return(&query.GetDayEventsResponse{Events: make([]event.Event, 0)}, nil)
				return a
			},
			request: &grpcapi.GetEventsRequest{
				Period: grpcapi.GetEventsRequest_GET_EVENTS_PERIOD_DAY,
			},
			wantErr: false,
			want:    &grpcapi.GetEventsResponse{Events: make([]*grpcapi.Event, 0)},
		},
	}
	for testName, tt := range tests {
		t.Run(testName, func(t *testing.T) {
			controller := gomock.NewController(t)
			defer controller.Finish()
			logger := mock.NewMockLogger(controller)
			s, err := internalgrpc.NewServer(tt.getApp(controller), logger)
			require.NoError(t, err)
			srv := s.(grpcapi.CalendarServer)

			getEventResponse, err := srv.GetEvents(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, getEventResponse)
				return
			}
			require.NoError(t, err)
			require.Equal(t, &grpcapi.GetEventsResponse{Events: make([]*grpcapi.Event, 0)}, getEventResponse)
		})
	}
}
