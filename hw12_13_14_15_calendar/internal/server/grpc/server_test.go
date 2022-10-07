package internalgrpc_test

import (
	"context"
	"testing"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/event/command"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/event/query"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/mock"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	internalgrpc "github.com/ekhvalov/hw12_13_14_15_calendar/internal/server/grpc"
	grpcapi "github.com/ekhvalov/hw12_13_14_15_calendar/pkg/api/grpc"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func Test_server_CreateEvent(t *testing.T) {
	eventTime := time.Now()
	tests := map[string]struct {
		getApp  func(controller *gomock.Controller) app.Application
		request *grpcapi.CreateEventRequest
		want    *grpcapi.CreateEventReply
		wantErr bool
	}{
		"when error occurred then should return error": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					CreateEvent(gomock.Any()).
					Return(nil, command.ErrValidate{})
				return a
			},
			wantErr: true,
		},
		"when no error occurred then should return CreateEventReply": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					CreateEvent(gomock.Any()).
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
			want: &grpcapi.CreateEventReply{
				Id:           "10",
				Title:        "Title",
				Date:         eventTime.Unix(),
				Duration:     30,
				Description:  "Description",
				NotifyBefore: 20,
			},
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

			createEventReply, err := srv.CreateEvent(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, createEventReply)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, createEventReply)
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
					UpdateEvent(gomock.Any()).
					Return(command.ErrValidate{})
				return a
			},
			wantErr: true,
		},
		"when no error occurred then should return Empty": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					UpdateEvent(gomock.Any()).
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

			updateEventReply, err := srv.UpdateEvent(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, updateEventReply)
				return
			}
			require.NoError(t, err)
			require.Equal(t, &grpcapi.Empty{}, updateEventReply)
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
					DeleteEvent(gomock.Any()).
					Return(event.ErrEventNotFound)
				return a
			},
			wantErr: true,
		},
		"when no error occurred then should return Empty": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					DeleteEvent(gomock.Any()).
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

			deleteEventReply, err := srv.DeleteEvent(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, deleteEventReply)
				return
			}
			require.NoError(t, err)
			require.Equal(t, &grpcapi.Empty{}, deleteEventReply)
		})
	}
}

func Test_server_GetEvents(t *testing.T) {
	tests := map[string]struct {
		getApp  func(controller *gomock.Controller) app.Application
		request *grpcapi.GetEventsRequest
		wantErr bool
		want    *grpcapi.GetEventsReply
	}{
		"when error occurred then should return error": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetDayEvents(gomock.Any()).
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
			want:    &grpcapi.GetEventsReply{Events: make([]*grpcapi.Event, 0)},
		},
		"when no error occurred then should return GetEventsReply": {
			getApp: func(controller *gomock.Controller) app.Application {
				a := mock.NewMockApplication(controller)
				a.EXPECT().
					GetDayEvents(gomock.Any()).
					Return(&query.GetDayEventsResponse{Events: make([]event.Event, 0)}, nil)
				return a
			},
			request: &grpcapi.GetEventsRequest{
				Period: grpcapi.GetEventsRequest_GET_EVENTS_PERIOD_DAY,
			},
			wantErr: false,
			want:    &grpcapi.GetEventsReply{Events: make([]*grpcapi.Event, 0)},
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

			getEventReply, err := srv.GetEvents(context.Background(), tt.request)

			if tt.wantErr {
				require.Error(t, err)
				require.Nil(t, getEventReply)
				return
			}
			require.NoError(t, err)
			require.Equal(t, &grpcapi.GetEventsReply{Events: make([]*grpcapi.Event, 0)}, getEventReply)
		})
	}
}
