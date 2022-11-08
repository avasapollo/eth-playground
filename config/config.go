package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NetworkURL string `envconfig:"NETWORK_URL"`
}

func Get() *Config {
	c := new(Config)
	_ = envconfig.Process("", c)
	return c
}
