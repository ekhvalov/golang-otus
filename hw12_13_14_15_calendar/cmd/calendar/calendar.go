package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/app"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/environment/config"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/server/http"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/storage"
	"github.com/spf13/cobra"
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
	v, err := config.CreateViper(cfgFile, configEnvPrefix, config.DefaultEnvKeyReplacer)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("create config error: %w", err))
	}
	cfg := NewConfig(v)

	logg := logger.New(cfg.Logger.Level, os.Stdout)

	strg, err := storage.CreateStorage(v)
	if err != nil {
		cobra.CheckErr(fmt.Errorf("create storage error: %w", err))
	}
	calendar, err := app.New(logg, strg)
	if err != nil {
		cobra.CheckErr(err)
	}

	httpServer := internalhttp.NewServer(calendar, logg)
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

	go func() {
		srvAddr := fmt.Sprintf("%s:%d", cfg.HTTP.Address, cfg.HTTP.Port)
		if err := httpServer.ListenAndServe(srvAddr); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
		}
		wg.Done()
	}()
	go func() {
		srvAddr := fmt.Sprintf("%s:%d", cfg.GRPC.Address, cfg.GRPC.Port)
		if err := grpcServer.ListenAndServe(srvAddr); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
		wg.Done()
	}()

	logg.Info("calendar is running...")
	wg.Wait()
}
