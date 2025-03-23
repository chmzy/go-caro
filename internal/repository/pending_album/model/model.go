package model

type AlbumPost struct {
	ID      int
	AlbumID string
	Author  string
	MsgLink Link
}

type Link struct {
	ChatID int64
	MsgID  int
}
