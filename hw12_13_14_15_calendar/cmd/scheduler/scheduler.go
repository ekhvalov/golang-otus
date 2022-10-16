package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app"
	appqueue "github.com/ekhvalov/hw12_13_14_15_calendar/internal/app/notification/queue"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/environment/config"
	configviper "github.com/ekhvalov/hw12_13_14_15_calendar/internal/environment/config/viper"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/environment/notification/queue/rabbitmq"
	memorystorage "github.com/ekhvalov/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"
)

var (
	configFile   string
	schedulerCmd = &cobra.Command{
		Use: "scheduler",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

const (
	configEnvPrefix = "scheduler"
	year            = time.Hour * 24 * 365
)

func init() {
	schedulerCmd.PersistentFlags().StringVar(&configFile, "config", "", "Path to config file")
}

func run() error {
	fmt.Println("Using config file:", configFile)
	configProvider, err := configviper.NewProvider(configFile, configEnvPrefix, configviper.DefaultEnvKeyReplacer)
	if err != nil {
		return fmt.Errorf("create config error: %w", err)
	}
	queueProducer, err := createQueueProducer(configProvider)
	if err != nil {
		return err
	}
	storage := createStorage()

	scheduler, err := app.NewScheduler(storage, queueProducer)
	if err != nil {
		return fmt.Errorf("create scheduler error: %w", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	errors := make(chan error, 2)
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		if err := scheduler.FindNotificationReadyEvents(ctx); err != nil {
			errors <- err
			cancel()
		}
		wg.Done()
	}()
	go func() {
		if err := scheduler.CleanOldEvents(ctx, year); err != nil {
			errors <- err
			cancel()
		}
		wg.Done()
	}()
	wg.Wait()
	close(errors)
	for e := range errors {
		err = multierror.Append(err, e)
	}
	return err
}

func createQueueProducer(provider config.Provider) (appqueue.Producer, error) {
	var cfg rabbitmq.ConfigRabbitMQ
	err := provider.UnmarshalKey("queue.rabbitmq", &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal ConfigRabbitMQ error: %w", err)
	}
	return rabbitmq.NewProducer(cfg), nil
}

func createStorage() event.Storage {
	return memorystorage.New(memorystorage.UUIDProvider{})
}
