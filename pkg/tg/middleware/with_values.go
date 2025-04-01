package middleware

import (
	"gopkg.in/telebot.v4"
)

func WithValues(values map[string]any, next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		for key, value := range values {
			c.Set(key, value)
		}
		return next(c)
	}
}
