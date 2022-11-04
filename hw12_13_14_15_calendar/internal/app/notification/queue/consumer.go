package queue

//go:generate mockgen -destination=./mock/consumer.gen.go -package mock . Consumer

import (
	"context"
	"io"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/notification"
)

// Consumer allows to subscribe to notification.Notification queue.
type Consumer interface {
	Subscribe(ctx context.Context) (<-chan notification.Notification, error)
	io.Closer
}
