package converter

import (
	"fmt"
	modelserv "go-caro/internal/service/history/model"
	modelapi "go-caro/pkg/tg/model"
)

func ToHistoryFromAPI(msg *modelapi.Message, ) *modelserv.PostHistory {
	return &modelserv.PostHistory{
		ID:       0,
		AlbumID:  msg.AlbumID,
		ChatID:   msg.Chat.ID,
		MsgID:    fmt.Sprintf("%d", msg.ID),
		PostedAt: msg.Time(),
	}
}
