package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/ekhvalov/hw12_13_14_15_calendar/internal/server/grpc"
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
	calendar, err := app.New(logg, storage)
	if err != nil {
		cobra.CheckErr(err)
	}

	httpServer := internalhttp.NewServer(logg, calendar)
	grpcServer, err := internalgrpc.NewServer(calendar, logg)
	if err != nil {
		cobra.CheckErr(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
		if err := grpcServer.Shutdown(ctx); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	// TODO: Run http and grpc
	go func(wg *sync.WaitGroup) {
		srvAddr := fmt.Sprintf("%s:%d", cfg.HTTP.Address, cfg.HTTP.Port)
		if err := httpServer.Start(srvAddr); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
		wg.Done()
	}(&wg)
	go func(wg *sync.WaitGroup) {
		srvAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Address, cfg.GRPC.Port)
		if err := grpcServer.Start(srvAddr); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
		wg.Done()
	}(&wg)

	logg.Info("calendar is running...")
	wg.Wait()
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
		return memorystorage.New(memorystorage.UUIDProvider{}), nil
	default:
		return nil, fmt.Errorf("undefined storage type: %s", cfg.Type)
	}
}
