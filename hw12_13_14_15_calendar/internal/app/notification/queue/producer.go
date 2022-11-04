package queue

//go:generate mockgen -destination=./mock/producer.gen.go -package mock . Producer

import "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/notification"

// ErrProducer will be returned in case of internal storage error.
type ErrProducer struct {
	message string
}

func (e ErrProducer) Error() string {
	return e.message
}

// NewErrProducer function creates ErrProducer error.
func NewErrProducer(message string) ErrProducer {
	return ErrProducer{message: message}
}

type Producer interface {
	Put(notification notification.Notification) error
	Close() error
}
