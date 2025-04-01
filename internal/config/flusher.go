package config

type FlushConfig interface {
	FlushPeriodSec() int64
}

type flushConfig struct {
	FlushPeriodSeconds int64 `yaml:"flush_period_sec"`
}

func NewFlusherConfig() (FlushConfig, error) {
	var cfg flushConfig

	if err := parseFile(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *flushConfig) FlushPeriodSec() int64 {
	return cfg.FlushPeriodSeconds
}
