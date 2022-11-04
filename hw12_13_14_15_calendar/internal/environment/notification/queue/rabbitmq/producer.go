package rabbitmq

import (
	"context"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/notification/queue"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/notification"
)

type producer struct {
	connector *connector
}

func NewProducer(mqConf ConfigRabbitMQ) queue.Producer {
	return &producer{connector: newConnector(mqConf.GetDSN(), mqConf.QueueName)}
}

func (p *producer) Put(notification notification.Notification) error {
	return p.connector.publish(context.Background(), notification)
}

func (p *producer) Close() error {
	return p.connector.close()
}
