package internalhttp

//nolint:lll
//go:generate oapi-codegen -package openapi -generate types -old-config-style -o ../../../pkg/api/openapi/types.gen.go ../../../api/openapi/openapi.yaml
//go:generate oapi-codegen -package openapi -generate spec -old-config-style -o ../../../pkg/api/openapi/spec.gen.go ../../../api/openapi/openapi.yaml
//go:generate oapi-codegen -package openapi -generate chi-server -old-config-style -o ../../../pkg/api/openapi/server.gen.go ../../../api/openapi/openapi.yaml
//go:generate oapi-codegen -package openapi -generate client -old-config-style -o ../../../pkg/api/openapi/client.gen.go ../../../api/openapi/openapi.yaml

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/server/http/event"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/pkg/api/openapi"
	"github.com/go-chi/chi/v5"
)

type Server interface {
	ListenAndServe(address string) error
	Shutdown(context context.Context) error
}

type server struct {
	logger app.Logger
	s      *http.Server
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

func NewServer(app app.Application, logger Logger) Server {
	eventsHandler := event.NewEventHandler(app, logger)

	router := chi.NewRouter()
	router.Use(loggingMiddleware(logger))

	openapi.HandlerFromMux(eventsHandler, router)

	return &server{
		logger: logger,
		s: &http.Server{
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

func (s *server) ListenAndServe(address string) error {
	s.logger.Info(fmt.Sprintf("listen: %s", address))
	s.s.Addr = address
	return s.s.ListenAndServe()
}

func (s *server) Shutdown(ctx context.Context) error {
	return s.s.Shutdown(ctx)
}
