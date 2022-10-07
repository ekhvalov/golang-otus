package main

import "github.com/spf13/viper"

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	HTTP    HTTPConf
	GRPC    GRPCConf
}

type LoggerConf struct {
	Level string
}

type StorageConf struct {
	Type string
}

type HTTPConf struct {
	Address string
	Port    uint
}

type GRPCConf struct {
	Address string
	Port    uint
}

func NewConfig(v *viper.Viper) Config {
	for key, value := range defaultValues {
		v.SetDefault(key, value)
	}
	return Config{
		Logger: LoggerConf{
			Level: v.GetString(loggerLevelKey),
		},
		Storage: StorageConf{
			Type: v.GetString(storageTypeKey),
		},
		HTTP: HTTPConf{
			Address: v.GetString(httpAddressKey),
			Port:    v.GetUint(httpPortKey),
		},
		GRPC: GRPCConf{
			Address: v.GetString(grpcAddressKey),
			Port:    v.GetUint(grpcPortKey),
		},
	}
}

const (
	grpcAddressKey = "grpc.address"
	grpcPortKey    = "grpc.port"
	httpAddressKey = "http.address"
	httpPortKey    = "http.port"
	loggerLevelKey = "logger.level"
	storageTypeKey = "storage.type"
)

var defaultValues = map[string]interface{}{
	loggerLevelKey: "info",
	storageTypeKey: "memory",
}
