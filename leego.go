package leego

import (
	"github.com/Leela0o5/LeeGo/config"
	"github.com/Leela0o5/LeeGo/engine"
	"github.com/Leela0o5/LeeGo/metrics"
)

type Config = config.Config

type Stats = metrics.Stats

func Run(cfg Config) *Stats {
	return engine.Run(cfg)
}

func RunAsync(cfg Config) (*Stats, chan struct{}) {
	return engine.RunAsync(cfg)
}
