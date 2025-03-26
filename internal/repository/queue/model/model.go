package model

type PostQueue struct {
	ID      int    `db:"id"`
	Author  string `db:"author"`
	AlbumID string `db:"album_id"`
	ChatID  int64  `db:"chat_id"`
	MsgID   string `db:"msg_id"`
}
