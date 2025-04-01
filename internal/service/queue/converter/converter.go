package converter

import (
	"fmt"
	modelserv "go-caro/internal/service/queue/model"
	modelapi "go-caro/pkg/tg/model"
)

func ToQueueFromAPI(msg *modelapi.Message) *modelserv.PostQueue {
	author := "unknown"
	if msg.Sender != nil {
		author = msg.Sender.FirstName
	}
	return &modelserv.PostQueue{
		ID:      0,
		Author:  author,
		AlbumID: msg.AlbumID,
		MsgID:   fmt.Sprintf("%d", msg.ID),
		ChatID:  msg.Chat.ID,
	}
}
