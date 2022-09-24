package internalhttp

import (
	"context"
	"fmt"
	"net/http"
)

type Server struct {
	srv    *http.Server
	logger Logger
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(address string, port uint, logger Logger, _ Application) *Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello world"))
	})

	middleware := loggingMiddleware(logger, handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", address, port),
		Handler: middleware,
	}

	return &Server{
		srv:    server,
		logger: logger,
	}
}

func (s *Server) Start(_ context.Context) error {
	s.logger.Info(fmt.Sprintf("listen: %s", s.srv.Addr))
	return s.srv.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
