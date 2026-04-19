package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	URL        string        `yaml:"url"`
	NumWorkers int           `yaml:"connections"`
	Duration   time.Duration `yaml:"duration"`
	RateLimit  int           `yaml:"rate"`
	Burst      int           `yaml:"burst"`
	Message    string        `yaml:"message"`
}

var defaults = Config{
	URL:        "ws://localhost:8080/ws",
	NumWorkers: 10,
	Duration:   10 * time.Second,
	RateLimit:  1000,
	Burst:      1000,
	Message:    "ping",
}

func LoadConfig(path string) (*Config, error) {
	cfg := defaults

	if path == "" {
		return &cfg, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("config: cannot open %q: %w", path, err)
	}
	defer f.Close()

	dec := yaml.NewDecoder(f)
	dec.KnownFields(true)

	if err := dec.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("config: failed to parse %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (c *Config) validate() error {
	if c.URL == "" {
		return fmt.Errorf("config: url is required")
	}
	if c.NumWorkers <= 0 {
		return fmt.Errorf("config: connections must be > 0, got %d", c.NumWorkers)
	}
	if c.Duration <= 0 {
		return fmt.Errorf("config: duration must be > 0, got %s", c.Duration)
	}
	if c.RateLimit < 0 {
		return fmt.Errorf("config: rate must be >= 0 (0 = unlimited), got %d", c.RateLimit)
	}
	return nil
}

func (c *Config) UnmarshalYAML(value *yaml.Node) error {
	type raw struct {
		URL        string `yaml:"url"`
		NumWorkers int    `yaml:"connections"`
		Duration   string `yaml:"duration"`
		RateLimit  int    `yaml:"rate"`
		Burst      int    `yaml:"burst"`
		Message    string `yaml:"message"`
	}

	var r raw
	if err := value.Decode(&r); err != nil {
		return err
	}

	c.URL = r.URL
	c.NumWorkers = r.NumWorkers
	c.RateLimit = r.RateLimit
	c.Burst = r.Burst
	c.Message = r.Message

	if r.Duration != "" {
		d, err := time.ParseDuration(r.Duration)
		if err != nil {
			return fmt.Errorf("config: invalid duration %q: %w", r.Duration, err)
		}
		c.Duration = d
	}

	return nil
}
