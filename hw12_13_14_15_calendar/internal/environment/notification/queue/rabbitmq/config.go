package rabbitmq

import (
	"fmt"
	"strings"
)

type ConfigRabbitMQ struct {
	Address   string
	Port      int
	Username  string
	Password  string
	QueueName string
}

func (c ConfigRabbitMQ) GetDSN() string {
	dsnBuilder := strings.Builder{}
	dsnBuilder.WriteString("amqp://")
	if c.Username != "" {
		dsnBuilder.WriteString(fmt.Sprintf("%s:%s@", c.Username, c.Password))
	}
	dsnBuilder.WriteString(fmt.Sprintf("%s:%d/", c.Address, c.Port))
	return dsnBuilder.String()
}
