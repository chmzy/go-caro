package tg

import (
	"fmt"
	"time"

	"gopkg.in/telebot.v4"
)

const (
	pollTimeout = 1 * time.Minute
)

type TgBot struct {
	Bot *telebot.Bot
}

func NewBot(token string) (*TgBot, error) {
	pref := telebot.Settings{
		Token:  token,
		Poller: &telebot.LongPoller{Timeout: pollTimeout},
	}

	// Create a new bot instance
	bot, err := telebot.NewBot(pref)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %v", err)
	}

	return &TgBot{
		Bot: bot,
	}, nil
}
