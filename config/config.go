package config

import (
	"time"
)

type Config struct {
	URL        string
	NumWorkers int
	Duration   time.Duration
	RateLimit int
}

func LoadConfig() (*Config, error) {
	return &Config{
		NumWorkers: 10,
		Duration:   10 * time.Second,
		RateLimit:  100,
	}, nil
}
