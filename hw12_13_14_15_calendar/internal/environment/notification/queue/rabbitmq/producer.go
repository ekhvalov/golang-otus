package rabbitmq

import (
	"context"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/notification/queue"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/notification"
)

type producer struct {
	connector *connector
}

func NewProducer(config Config) queue.Producer {
	return &producer{connector: newConnector(config.GetDSN(), config.QueueName)}
}

func (p *producer) Put(notification notification.Notification) error {
	return p.connector.publish(context.Background(), notification)
}

func (p *producer) Close() error {
	return p.connector.close()
}
