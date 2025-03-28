package middleware

import (
	"gopkg.in/telebot.v4"
)

func ForwardedFromChannel(chan_id int64, next, fallback telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if !c.Message().IsForwarded() {
			return fallback(c)
		}

		if c.Message().OriginalChat == nil {
			return fallback(c)
		}

		if c.Message().OriginalChat.ID != chan_id {
			return fallback(c)
		}

		return next(c)
	}
}
