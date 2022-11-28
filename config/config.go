package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NetworkURL          string `envconfig:"NETWORK_URL"`
	MyAccount           string `envconfig:"MY_ACCOUNT"`
	MyAccountPrivateKey string `envconfig:"MY_ACCOUNT_PRIVATE_KEY"`

	KeyStoreAccountPassword string `envconfig:"KEY_STORE_ACCOUNT_PASSWORD"`
	KeyStorePath            string `envconfig:"KEY_STORE_PATH"`
}

func Get() *Config {
	c := new(Config)
	_ = envconfig.Process("", c)
	return c
}
