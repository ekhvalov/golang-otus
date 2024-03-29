package rabbitmq

import (
	"context"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/notification/queue"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/notification"
)

func NewConsumer(config Config) queue.Consumer {
	return &consumer{connector: newConnector(config.GetDSN(), config.QueueName)}
}

type consumer struct {
	connector *connector
}

func (c *consumer) Subscribe(ctx context.Context) (<-chan notification.Notification, error) {
	return c.connector.consume(ctx)
}

func (c *consumer) Close() error {
	if c.connector != nil {
		return c.connector.close()
	}
	return nil
}
