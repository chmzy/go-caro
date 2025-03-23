package model

type PostQueue struct {
	ID      int
	Author  string
	AlbumID string
	MsgLink Link
}

type Link struct {
	ChatID int64
	MsgID  string
}
