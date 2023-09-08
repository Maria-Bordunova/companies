package config

import (
	"github.com/jinzhu/configor"
)

const appEnvPrefix = "C"

func InitConfig() (*Config, error) {
	var cfg = &Config{}

	err := configor.
		New(&configor.Config{ENVPrefix: appEnvPrefix}).Load(cfg)

	if err != nil {
		return nil, err
	}

	return cfg, nil
}
