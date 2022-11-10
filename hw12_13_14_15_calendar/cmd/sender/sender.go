package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/environment/config"
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
	v, err := config.CreateViper(configFile, configEnvPrefix, config.DefaultEnvKeyReplacer)
	if err != nil {
		return fmt.Errorf("create viper error: %w", err)
	}
	queueConsumer := rabbitmq.NewConsumer(config.NewRabbitMQConfig(v))
	defer func() {
		err = queueConsumer.Close()
	}()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	notifications, err := queueConsumer.Subscribe(ctx)
	if err != nil {
		return fmt.Errorf("consumer subscribe error: %w", err)
	}
	c := config.NewSenderWriterConfig(v)
	var output io.Writer
	if c.TargetFile == "" {
		output = os.Stdout
	} else {
		file, err := os.Create(c.TargetFile)
		if err != nil {
			return fmt.Errorf("create file '%s' error: %w", c.TargetFile, err)
		}
		defer func() {
			err = file.Close()
		}()
		output = file
	}
	s := sender.NewSender(output)
	for notification := range notifications {
		if err = s.Send(notification); err != nil {
			return fmt.Errorf("send error: %w", err)
		}
	}
	return err
}
