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
	output := &cloneWriter{}
	output.addWriter(os.Stdout)

	if c.TargetFile != "" {
		file, err := os.Create(c.TargetFile)
		if err != nil {
			return fmt.Errorf("create file '%s' error: %w", c.TargetFile, err)
		}
		defer func() {
			err = file.Close()
		}()
		output.addWriter(file)
	}
	s := sender.NewSender(output)
	for notification := range notifications {
		if err = s.Send(notification); err != nil {
			return fmt.Errorf("send error: %w", err)
		}
	}
	return err
}

type cloneWriter struct {
	writers []io.Writer
}

func (w *cloneWriter) Write(p []byte) (n int, err error) {
	for _, writer := range w.writers {
		n, err = writer.Write(p)
		if err != nil {
			return n, fmt.Errorf("write error: %w", err)
		}
	}
	return n, nil
}

func (w *cloneWriter) addWriter(writer io.Writer) {
	if w.writers == nil {
		w.writers = make([]io.Writer, 0)
	}
	w.writers = append(w.writers, writer)
}
