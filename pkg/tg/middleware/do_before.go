package middleware

import (
	"gopkg.in/telebot.v4"
)

func DoBefore(action func(), next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		action()
		return next(c)
	}
}
