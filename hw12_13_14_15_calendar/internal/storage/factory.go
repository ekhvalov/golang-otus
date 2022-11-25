package storage

import (
	"fmt"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/domain/event"
	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/environment/config"
	memorystorage "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/storage/memory"
	pgsqlstorage "github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/storage/pgsql"
	"github.com/spf13/viper"
)

func CreateStorage(v *viper.Viper) (event.Storage, error) {
	sc := config.NewStorageConfig(v)
	storageType, err := sc.GetStorageType()
	if err != nil {
		return nil, err
	}
	switch storageType {
	case config.MEMORY:
		return memorystorage.New(memorystorage.UUIDProvider{}), nil
	case config.PGSQL:
		conf := config.CreatePgsqlConfig(v)
		return pgsqlstorage.NewStorage(conf), nil
	default:
		return nil, fmt.Errorf("undefined storage type: %s", storageType)
	}
}
