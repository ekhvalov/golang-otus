package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

var DefaultEnvKeyReplacer *strings.Replacer

func init() {
	DefaultEnvKeyReplacer = strings.NewReplacer(".", "_")
}

func CreateViper(configFile, envPrefix string, envKeyReplacer *strings.Replacer) (*viper.Viper, error) {
	v := viper.New()
	if envPrefix != "" {
		v.SetEnvPrefix(envPrefix)
	}
	if envKeyReplacer != nil {
		v.SetEnvKeyReplacer(envKeyReplacer)
	}
	v.AutomaticEnv()

	if configFile == "" {
		err := v.ReadConfig(bytes.NewBuffer([]byte("")))
		if err != nil {
			return nil, err
		}
	} else {
		v.SetConfigFile(configFile)
		err := v.ReadInConfig()
		if err != nil {
			return nil, fmt.Errorf("config file '%s' error: %w", configFile, err)
		}
		fmt.Println("Using config file:", configFile)
	}
	return v, nil
}
