package config

import "time"

type AuthConfig struct {
	SecretKey string
	Duration  time.Duration
}
