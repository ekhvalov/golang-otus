package config

import (
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/environment/notification/queue/rabbitmq"
	"github.com/spf13/viper"
)

func NewRabbitMQConfig(v *viper.Viper) rabbitmq.Config {
	return rabbitmq.Config{
		Address:   v.GetString("queue.rabbitmq.address"),
		Port:      v.GetInt("queue.rabbitmq.port"),
		Username:  v.GetString("queue.rabbitmq.username"),
		Password:  v.GetString("queue.rabbitmq.password"),
		QueueName: v.GetString("queue.rabbitmq.queue_name"),
	}
}
