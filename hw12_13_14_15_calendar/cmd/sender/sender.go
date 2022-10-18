package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/environment/config"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/environment/notification/queue/rabbitmq"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/environment/notification/sender"
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
	v, err := config.CreateViper(configFile, configEnvPrefix, config.DefaultEnvKeyReplacer)
	if err != nil {
		return fmt.Errorf("create viper error: %w", err)
	}
	queueConsumer := rabbitmq.NewConsumer(config.NewRabbitMQConfig(v))
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
