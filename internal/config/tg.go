package config

type TgConfig interface {
	Token() string
	ID() int64
	Admins() []string
	PostTimeout() int64
	RepostTimeout() int64
}

type tgConfig struct {
	BotToken          string   `yaml:"token"`
	ChannelID         int64    `yaml:"channel_id"`
	AdminUsers        []string `yaml:"admins"`
	RepostTimeoutHour int64    `yaml:"repost_timeout"`
	PostTimeoutHour   int64    `yaml:"post_timeout"`
}

func NewTGConfig() (TgConfig, error) {
	var cfg tgConfig

	if err := parseFile(configPath, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *tgConfig) ID() int64 {
	return cfg.ChannelID
}

func (cfg *tgConfig) Token() string {
	return cfg.BotToken
}

func (cfg *tgConfig) Admins() []string {
	return cfg.AdminUsers
}

func (cfg *tgConfig) RepostTimeout() int64 {
	return cfg.RepostTimeoutHour
}
func (cfg *tgConfig) PostTimeout() int64 {
	return cfg.PostTimeoutHour
}
