package config

type TgConfig interface {
	Token() string
	ChannelID() int64
	SuggestionsID() int64
	Admins() []string
	PostPeriodSec() int64
	RepostPeriodSec() int64
}

type tgConfig struct {
	BotToken             string   `yaml:"token"`
	MainChannelID        int64    `yaml:"channel_id"`
	SuggestionsChannelID int64    `yaml:"suggestions_id"`
	AdminUsers           []string `yaml:"admins"`
	RepostPeriodSecond   int64    `yaml:"repost_period_sec"`
	PostPeriodSecond     int64    `yaml:"post_period_sec"`
}

func NewTGConfig() (TgConfig, error) {
	var cfg tgConfig

	if err := parseFile(configPath, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *tgConfig) ChannelID() int64 {
	return cfg.MainChannelID
}

func (cfg *tgConfig) SuggestionsID() int64 {
	return cfg.SuggestionsChannelID
}
func (cfg *tgConfig) Token() string {
	return cfg.BotToken
}

func (cfg *tgConfig) Admins() []string {
	return cfg.AdminUsers
}

func (cfg *tgConfig) RepostPeriodSec() int64 {
	return cfg.RepostPeriodSecond
}
func (cfg *tgConfig) PostPeriodSec() int64 {
	return cfg.PostPeriodSecond
}
