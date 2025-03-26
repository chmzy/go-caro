package model

import "time"

type PostHistory struct {
	ID       int
	AlbumID  string
	PostedAt time.Time
}
