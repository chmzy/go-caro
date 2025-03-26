package model

import "gopkg.in/telebot.v4"

type Message = telebot.Message
type Post = telebot.StoredMessage
type Context = telebot.Context

var NOOP = func(ctx telebot.Context) error { return nil }
