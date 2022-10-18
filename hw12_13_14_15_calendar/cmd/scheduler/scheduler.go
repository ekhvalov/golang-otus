package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/environment/config"
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
	v, err := config.CreateViper(configFile, configEnvPrefix, config.DefaultEnvKeyReplacer)
	if err != nil {
		return fmt.Errorf("create viper error: %w", err)
	}
	queueProducer := rabbitmq.NewProducer(config.NewRabbitMQConfig(v))
	storage := memorystorage.New(memorystorage.UUIDProvider{})

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
