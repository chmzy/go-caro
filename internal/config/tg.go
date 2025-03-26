package config

type TgConfig interface {
	Token() string
	ChannelID() int64
	SuggestionsID() int64
	Admins() []string
	PostTimeout() int64
	RepostTimeout() int64
}

type tgConfig struct {
	BotToken             string   `yaml:"token"`
	MainChannelID        int64    `yaml:"channel_id"`
	SuggestionsChannelID int64    `yaml:"suggestions_id"`
	AdminUsers           []string `yaml:"admins"`
	RepostTimeoutHour    int64    `yaml:"repost_timeout"`
	PostTimeoutHour      int64    `yaml:"post_timeout"`
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

func (cfg *tgConfig) RepostTimeout() int64 {
	return cfg.RepostTimeoutHour
}
func (cfg *tgConfig) PostTimeout() int64 {
	return cfg.PostTimeoutHour
}
