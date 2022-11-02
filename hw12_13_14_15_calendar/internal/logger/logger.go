package logger

import (
	"io"
	"strings"

	"github.com/rs/zerolog"
)

const (
	levelDebug   = "debug"
	levelInfo    = "info"
	levelWarning = "warning"
	levelError   = "error"
)

type Logger struct {
	logger zerolog.Logger
}

func New(level string, w io.Writer) *Logger {
	l := zerolog.New(w).With().Timestamp().Logger()
	switch strings.ToLower(level) {
	case levelDebug:
		l = l.Level(zerolog.DebugLevel)
	case "warn":
		fallthrough
	case levelWarning:
		l = l.Level(zerolog.WarnLevel)
	case levelError:
		l = l.Level(zerolog.ErrorLevel)
	default:
		l = l.Level(zerolog.InfoLevel)
	}
	return &Logger{logger: l}
}

func (l Logger) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l Logger) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error().Msg(msg)
}
