package middleware

import (
	"gopkg.in/telebot.v4"
)

func ForwardedFromChannel(chan_id int64, next, fallback telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.Message().IsForwarded() && c.Message().OriginalChat.ID == chan_id {
			return next(c)
		}

		return fallback(c)
	}
}
