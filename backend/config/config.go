package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port            string        `envconfig:"PORT" required:"true"`
	DbUrl           string        `envconfig:"DBURL" required:"true"`
	DbName          string        `envconfig:"DBNAME" required:"true"`
	DbUsername      string        `envconfig:"DBUSERNAME"`
	DbPassword      string        `envconfig:"DBPASSWORD"`
	SecretKey       string        `envconfig:"SECRETKEY" required:"true"`
	Duration        time.Duration `envconfig:"DURATION" default:"1h"`
	Origins         []string      `envconfig:"WEB_ORIGIN"`
	UserIconsBucket string        `envconfig:"USERICONS_BUCKET" required:"true"`
}

const EnvPrefix = ""

func NewConfig() (Config, error) {
	var cfg Config

	if err := envconfig.Process(EnvPrefix, &cfg); err != nil {
		return Config{}, fmt.Errorf("procces new config: %w", err)
	}
	return cfg, nil
}

func (c Config) StoreConfig() StoreConfig {
	return StoreConfig{
		DbUrl:      c.DbUrl,
		DbName:     c.DbName,
		DbUsername: c.DbUsername,
		DbPassword: c.DbPassword,
	}
}

func (c Config) ServerConfig() ServerConfig {
	return ServerConfig{
		Port:    c.Port,
		Origins: c.Origins,
	}
}

func (c Config) AuthConfig() AuthConfig {
	return AuthConfig{
		SecretKey: c.SecretKey,
		Duration:  c.Duration,
	}
}

func (c Config) BucketConfig() BucketConfig {
	return BucketConfig{
		UserIconsBucket: c.UserIconsBucket,
	}
}
