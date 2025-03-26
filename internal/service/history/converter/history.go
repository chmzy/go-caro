package converter

import (
	modelserv "go-caro/internal/service/history/model"
	modelapi "go-caro/pkg/tg/model"
)

func ToHistoryFromAPI(msg *modelapi.Message) *modelserv.PostHistory {
	return &modelserv.PostHistory{
		ID:       msg.ID,
		AlbumID:  msg.AlbumID,
		PostedAt: msg.Time(),
	}
}
