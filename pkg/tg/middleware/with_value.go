package middleware

import (
	"gopkg.in/telebot.v4"
)

func WithValue(key string, value interface{}, next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		c.Set(key, value)
		return next(c)
	}
}
