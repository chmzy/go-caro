package config

type PGConfig interface {
	DSN() string
	Timezone() int64
}

type pgConfig struct {
	Dsn         string `yaml:"dsn"`
	GMTTimezone int64  `yaml:"gmt_timezone"`
}

func NewPGConfig() (PGConfig, error) {
	var cfg pgConfig

	if err := parseFile(configPath, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *pgConfig) DSN() string {
	return cfg.Dsn
}
func (cfg *pgConfig) Timezone() int64 {
	return cfg.GMTTimezone
}
