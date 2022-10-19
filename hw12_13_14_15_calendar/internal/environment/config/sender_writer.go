package config

import (
	"github.com/ekhvalov/hw12_13_14_15_calendar/internal/environment/notification/sender"
	"github.com/spf13/viper"
)

func NewSenderWriterConfig(v *viper.Viper) sender.Config {
	return sender.Config{TargetFile: v.GetString("writer.target_file")}
}
