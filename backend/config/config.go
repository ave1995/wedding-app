package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port   string `envconfig:"PORT" required:"true"`
	DbUrl  string `envconfig:"DBURL" required:"true"`
	DbName string `envconfig:"DBNAME" required:"true"`
}

const ENV_PREFIX = "backend"

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process(ENV_PREFIX, &config); err != nil {
		return Config{}, fmt.Errorf("procces new config: %w", err)
	}
	return config, nil
}

func (c Config) StoreConfig() StoreConfig {
	return StoreConfig{
		DbUrl:  c.DbUrl,
		DbName: c.DbName,
	}
}

func (c Config) ServerConfig() ServerConfig {
	return ServerConfig{
		Port: c.Port,
	}
}
