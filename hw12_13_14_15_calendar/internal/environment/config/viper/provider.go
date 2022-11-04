package configviper

import (
	"fmt"
	"strings"

	"github.com/ekhvalov/golang-otus/hw12_13_14_15_calendar/internal/environment/config"
	"github.com/spf13/viper"
)

var DefaultEnvKeyReplacer *strings.Replacer

func init() {
	DefaultEnvKeyReplacer = strings.NewReplacer(".", "_")
}

func NewProvider(configFile, envPrefix string, envKeyReplacer *strings.Replacer) (config.Provider, error) {
	v, err := createViper(configFile, envPrefix, envKeyReplacer)
	if err != nil {
		return nil, fmt.Errorf("create viper error: %w", err)
	}
	return &provider{v: v}, nil
}

type provider struct {
	v *viper.Viper
}

func (p *provider) UnmarshalKey(key string, rawVal interface{}) error {
	return p.v.UnmarshalKey(key, rawVal)
}

func (p *provider) GetString(key string) string {
	return p.v.GetString(key)
}

func createViper(configFile, envPrefix string, envKeyReplacer *strings.Replacer) (*viper.Viper, error) {
	v := viper.New()
	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
	}
	if envKeyReplacer != nil {
		v.SetEnvKeyReplacer(envKeyReplacer)
	}
	v.AutomaticEnv()
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("config file '%s' error: %w", configFile, err)
	}
	return v, nil
}
