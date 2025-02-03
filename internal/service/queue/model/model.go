package model

type PostQueue struct {
	ID      int
	MsgLink Link
}

type Link struct {
	ChatID int64
	MsgID  string
}
