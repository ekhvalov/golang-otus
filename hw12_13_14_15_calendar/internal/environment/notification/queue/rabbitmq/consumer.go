package rabbitmq

import (
	"context"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/notification/queue"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/notification"
)

func NewConsumer(conf ConfigRabbitMQ) queue.Consumer {
	return &consumer{connector: newConnector(conf.GetDSN(), conf.QueueName)}
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
