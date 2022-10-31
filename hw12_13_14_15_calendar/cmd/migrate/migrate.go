package main

import (
	"fmt"

	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/environment/config"
	pgsqlstorage "github.com/ekhvalov/hw12_13_14_15_calendar/internal/storage/pgsql"
	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	command    string
	migrateCmd = &cobra.Command{
		Use: "migrate",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

const configEnvPrefix = "calendar"

func init() {
	migrateCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Path to config file")
	migrateCmd.PersistentFlags().StringVar(&command, "command", "", "Migrate command (up, down, etc)")
}

func run() error {
	v, err := config.CreateViper(cfgFile, configEnvPrefix, config.DefaultEnvKeyReplacer)
	if err != nil {
		return fmt.Errorf("viper config create error: %w", err)
	}
	storageConfig := config.NewStorageConfig(v)
	storageType, err := storageConfig.GetStorageType()
	if err != nil {
		return fmt.Errorf("get storage type error: %w", err)
	}
	switch storageType {
	case config.MEMORY:
		return nil
	case config.PGSQL:
		conf := config.CreatePgsqlConfig(v)
		m := pgsqlstorage.NewMigrator(conf)
		return m.Run(command)
	default:
		return fmt.Errorf("unknown storage type: %s", storageType)
	}
}
