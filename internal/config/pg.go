package config

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	Dsn string `yaml:"dsn"`
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
