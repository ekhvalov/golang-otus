package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	appqueue "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app/notification/queue"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/environment/config"
	configviper "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/environment/config/viper"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/environment/notification/queue/rabbitmq"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/environment/notification/sender"
	"github.com/spf13/cobra"
)

var (
	configFile string
	senderCmd  = &cobra.Command{
		Use: "sender",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

const configEnvPrefix = "sender"

func init() {
	senderCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to config file")
}

func run() error {
	fmt.Println("Using config file:", configFile)
	configProvider, err := configviper.NewProvider(configFile, configEnvPrefix, configviper.DefaultEnvKeyReplacer)
	if err != nil {
		return fmt.Errorf("create config error: %w", err)
	}
	queueConsumer, err := createQueueConsumer(configProvider)
	if err != nil {
		return fmt.Errorf("create consumer error: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	notifications, err := queueConsumer.Subscribe(ctx)
	if err != nil {
		return fmt.Errorf("consumer subscribe error: %w", err)
	}
	s := sender.NewSender(os.Stdout)
	for notification := range notifications {
		if err = s.Send(notification); err != nil {
			return fmt.Errorf("send error: %w", err)
		}
	}
	return nil
}

func createQueueConsumer(provider config.Provider) (appqueue.Consumer, error) {
	var cfg rabbitmq.ConfigRabbitMQ
	err := provider.UnmarshalKey("queue.rabbitmq", cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal ConfigRabbitMQ error: %w", err)
	}
	return rabbitmq.NewConsumer(cfg), nil
}
