package internalgrpc

//nolint:lll
//go:generate protoc ../../../api/grpc/EventService.proto -I ../../../api/grpc --go_out=../../../pkg/api/grpc --go-grpc_out=../../../pkg/api/grpc

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/command"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/event/query"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	grpcapi "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/pkg/api/grpc"
	"google.golang.org/grpc"
)

type Server interface {
	Start(address string) error
	Shutdown(context context.Context) error
}

func NewServer(application app.Application, logger app.Logger) (Server, error) {
	if application == nil {
		return nil, fmt.Errorf("expect app.Application <nil> provided")
	}
	if logger == nil {
		return nil, fmt.Errorf("expect app.Logger <nil> provided")
	}
	return &server{
		application: application,
		logger:      logger,
	}, nil
}

type server struct {
	grpcapi.UnimplementedCalendarServer
	application app.Application
	logger      app.Logger
	grpcServer  *grpc.Server
}

func (s *server) CreateEvent(
	_ context.Context,
	request *grpcapi.CreateEventRequest,
) (*grpcapi.CreateEventReply, error) {
	createEventRequest := command.CreateEventRequest{
		Title:        request.GetTitle(),
		DateTime:     time.Unix(request.GetDate(), 0),
		Duration:     time.Duration(request.GetDuration()) * time.Minute,
		UserID:       "10",
		Description:  request.GetDescription(),
		NotifyBefore: time.Duration(request.GetNotifyBefore()) * time.Minute,
	}
	newEvent, err := s.application.CreateEvent(createEventRequest)
	if err != nil {
		return nil, fmt.Errorf("create event error: %w", err)
	}
	return &grpcapi.CreateEventReply{
		Id:           newEvent.Event.ID,
		Title:        newEvent.Event.Title,
		Date:         newEvent.Event.DateTime.Unix(),
		Duration:     int32(newEvent.Event.Duration / time.Minute),
		Description:  newEvent.Event.Description,
		NotifyBefore: int32(newEvent.Event.NotifyBefore / time.Minute),
	}, nil
}

func (s *server) UpdateEvent(_ context.Context, request *grpcapi.UpdateEventRequest) (*grpcapi.Empty, error) {
	updateEventRequest := command.UpdateEventRequest{
		ID:           request.GetId(),
		Title:        request.GetTitle(),
		DateTime:     time.Unix(request.GetDate(), 0),
		Duration:     time.Duration(request.GetDuration()) * time.Minute,
		UserID:       "10",
		Description:  request.GetDescription(),
		NotifyBefore: time.Duration(request.GetNotifyBefore()) * time.Minute,
	}
	if err := s.application.UpdateEvent(updateEventRequest); err != nil {
		return nil, fmt.Errorf("update event error: %w", err)
	}
	return &grpcapi.Empty{}, nil
}

func (s *server) DeleteEvent(_ context.Context, request *grpcapi.DeleteEventRequest) (*grpcapi.Empty, error) {
	deleteEventRequest := command.DeleteEventRequest{ID: request.GetId()}
	if err := s.application.DeleteEvent(deleteEventRequest); err != nil {
		return nil, fmt.Errorf("delete event error: %w", err)
	}
	return &grpcapi.Empty{}, nil
}

func (s *server) GetEvents(_ context.Context, request *grpcapi.GetEventsRequest) (*grpcapi.GetEventsReply, error) {
	var events []event.Event
	switch request.GetPeriod() {
	case grpcapi.GetEventsRequest_GET_EVENTS_PERIOD_UNSPECIFIED:
		return nil, fmt.Errorf("get events period is unspecified")
	case grpcapi.GetEventsRequest_GET_EVENTS_PERIOD_DAY:
		response, err := s.application.GetDayEvents(query.GetDayEventsRequest{Date: time.Unix(request.GetDate(), 0)})
		if err != nil {
			return nil, fmt.Errorf("get day events error: %w", err)
		}
		events = response.Events
	case grpcapi.GetEventsRequest_GET_EVENTS_PERIOD_WEEK:
		response, err := s.application.GetWeekEvents(query.GetWeekEventsRequest{Date: time.Unix(request.GetDate(), 0)})
		if err != nil {
			return nil, fmt.Errorf("get week events error: %w", err)
		}
		events = response.Events
	case grpcapi.GetEventsRequest_GET_EVENTS_PERIOD_MONTH:
		response, err := s.application.GetMonthEvents(query.GetMonthEventsRequest{Date: time.Unix(request.GetDate(), 0)})
		if err != nil {
			return nil, fmt.Errorf("get month events error: %w", err)
		}
		events = response.Events
	}
	e := make([]*grpcapi.Event, len(events))
	for i, ev := range events {
		e[i] = &grpcapi.Event{
			Id:           ev.ID,
			Title:        ev.Title,
			Date:         ev.DateTime.Unix(),
			Duration:     int32(ev.Duration / time.Minute),
			Description:  ev.Description,
			NotifyBefore: int32(ev.NotifyBefore / time.Minute),
		}
	}
	return &grpcapi.GetEventsReply{Events: e}, nil
}

func (s *server) Start(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return fmt.Errorf("listen address '%s' error: %w", address, err)
	}
	s.grpcServer = grpc.NewServer()
	grpcapi.RegisterCalendarServer(s.grpcServer, s)
	s.logger.Info(fmt.Sprintf("listen on: %s", address))
	return s.grpcServer.Serve(listener)
}

func (s *server) Shutdown(_ context.Context) error {
	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
		s.logger.Info("stopped successfully")
	}
	return nil
}
