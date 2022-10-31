package config

import (
	"fmt"
	"strings"

	pgsqlstorage "github.com/ekhvalov/hw12_13_14_15_calendar/internal/storage/pgsql"
	"github.com/spf13/viper"
)

type StorageType string

const (
	MEMORY StorageType = "memory"
	PGSQL  StorageType = "pgsql"
)

func NewStorageConfig(v *viper.Viper) *StorageConfig {
	return &StorageConfig{v: v}
}

type StorageConfig struct {
	v *viper.Viper
}

func (c StorageConfig) GetStorageType() (StorageType, error) {
	storageType := c.v.GetString("storage.type")
	switch StorageType(strings.ToLower(storageType)) {
	case MEMORY:
		return MEMORY, nil
	case PGSQL:
		return PGSQL, nil
	default:
		return "", fmt.Errorf("unknown storage type: %s", storageType)
	}
}

func CreatePgsqlConfig(v *viper.Viper) pgsqlstorage.Config {
	return pgsqlstorage.Config{
		Host:     v.GetString("storage.pgsql.host"),
		Port:     v.GetInt("storage.pgsql.port"),
		Username: v.GetString("storage.pgsql.username"),
		Password: v.GetString("storage.pgsql.password"),
	}
}
