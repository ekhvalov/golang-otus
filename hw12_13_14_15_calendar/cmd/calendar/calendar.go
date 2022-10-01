package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ekhvalov/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/ekhvalov/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	calendarCmd = &cobra.Command{
		Use: "calendar",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
)

const configEnvPrefix = "calendar"

func init() {
	calendarCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to config file")
}

func run() {
	v, err := getViper()
	if err != nil {
		cobra.CheckErr(fmt.Errorf("create config error: %w", err))
	}
	cfg := NewConfig(v)

	logg := logger.New(cfg.Logger.Level, os.Stdout)

	storage, err := createStorage(cfg.Storage)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("create storage error: %w", err))
	}
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(cfg.HTTP.Address, cfg.HTTP.Port, logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

func getViper() (*viper.Viper, error) {
	v := viper.New()
	v.SetEnvPrefix(configEnvPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if cfgFile == "" {
		err := v.ReadConfig(bytes.NewBuffer([]byte("")))
		if err != nil {
			return nil, err
		}
	} else {
		v.SetConfigFile(cfgFile)
		err := v.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("config file '%s' error: %w", cfgFile, err)
		}
		fmt.Println("Using config file:", cfgFile)
	}
	return v, nil
}

func createStorage(cfg StorageConf) (event.Storage, error) {
	switch strings.ToLower(cfg.Type) {
	case "memory":
		return memorystorage.New(), nil
	default:
		return nil, fmt.Errorf("undefined storage type: %s", cfg.Type)
	}
}
