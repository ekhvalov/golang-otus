//go:build integration
// +build integration

package integration_test

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/notification"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/require"
)

const (
	exchangeName        = ""
	contentType         = "application/octet-stream"
	defaultRabbitMQHost = "localhost"
	defaultRabbitMQPort = "5672"
	defaultQueueName    = "events_notifications"
	defaultTargetFile   = "/tmp/writer.txt"
)

func Test_SendNotification(t *testing.T) {
	time.Sleep(time.Second * 10)

	connection, err := amqp.Dial(getAMQPDsn())
	require.NoError(t, err)
	ch, err := connection.Channel()
	require.NoError(t, err)
	buffer := bytes.Buffer{}
	encoder := gob.NewEncoder(&buffer)
	n := notification.Notification{
		EventID:    "100500",
		EventTitle: "Event 1",
		EventDate:  time.Now(),
		UserID:     "100600",
	}
	err = encoder.Encode(n)
	require.NoError(t, err)

	err = ch.PublishWithContext(
		context.Background(),
		exchangeName,
		getEnv("TESTS_QUEUE_NAME", defaultQueueName),
		false,
		false,
		amqp.Publishing{
			ContentType: contentType,
			Body:        buffer.Bytes(),
		})
	require.NoError(t, err)

	time.Sleep(time.Second * 3)

	data, err := os.ReadFile(getEnv("TESTS_WRITER_TARGET_FILE", defaultTargetFile))
	require.NoError(t, err)

	require.Contains(t, string(data), n.EventTitle)
}

func getEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		return defaultValue
	}
	return value
}

func getAMQPDsn() string {
	host := getEnv("TESTS_RQBBITMQ_HOST", defaultRabbitMQHost)
	port := getEnv("TESTS_RABBITMQ_PORT", defaultRabbitMQPort)
	return fmt.Sprintf("amqp://guest:guest@%s:%s", host, port)
}
