package middleware

import (
	"gopkg.in/telebot.v4"
)

func IsForwarded(next, fallback telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		if c.Message().IsForwarded() {
			return next(c)
		}

		return fallback(c)
	}
}
