package middleware

import (
	"gopkg.in/telebot.v4"
)

func FromAdmin(admins []string, next, fallback telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		for _, admin := range admins {
			if c.Sender().Username == admin {
				return next(c)
			}
		}

		return fallback(c)
	}
}
