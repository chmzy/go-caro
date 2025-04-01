package model

import "time"

type PostHistory struct {
	ID       int
	AlbumID  string
	ChatID   int64
	MsgID    string
	PostedAt time.Time
}
