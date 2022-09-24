package main

import "github.com/spf13/viper"

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

type Config struct {
	Logger LoggerConf
	// TODO
}

type LoggerConf struct {
	Level string
	// TODO
}

func NewConfig(v *viper.Viper) Config {
	for key, value := range defaultValues {
		v.SetDefault(key, value)
	}
	return Config{
		Logger: LoggerConf{
			Level: v.GetString(loggerLevel),
		},
	}
}

const (
	loggerLevel = "logger.level"
)

var defaultValues = map[string]interface{}{
	loggerLevel: "info",
}
