package converter

import (
	"fmt"
	modelserv "go-caro/internal/service/queue/model"
	modelapi "go-caro/pkg/tg/model"
)

func ToQueueFromAPI(msg *modelapi.Message) *modelserv.PostQueue {
	return &modelserv.PostQueue{
		ID: 0,
		MsgLink: modelserv.Link{
			MsgID:  fmt.Sprint(msg.ID),
			ChatID: msg.Chat.ID,
		},
	}
}
