package config

import (
	"github.com/caarlos0/env/v6"
)

type Config struct {
	MemcachedEndpoint string `env:"MEMCACHED_ENDPOINT" envDefault:"localhost:11211"`
	Port              int    `env:"PORT" envDefault:"44444"`
}

func New() (*Config, error) {
	cfg := new(Config)
	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
