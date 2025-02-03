package config

type FlushConfig interface {
	FlushTimeout() int64
}

type flushConfig struct {
	FlushTimeoutHour int64 `yaml:"flush_timeout"`
}

func NewFlusherConfig() (FlushConfig, error) {
	var cfg flushConfig

	if err := parseFile(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *flushConfig) FlushTimeout() int64 {
	return cfg.FlushTimeoutHour
}
