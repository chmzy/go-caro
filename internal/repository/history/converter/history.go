package converter

import (
	modelrepo "go-caro/internal/repository/history/model"
	modelserv "go-caro/internal/service/history/model"
)

func ToHistoryFromRepo(post *modelrepo.PostHistory) *modelserv.PostHistory {
	return &modelserv.PostHistory{
		ID:       post.ID,
		AlbumID: post.AlbumID,
		PostedAt: post.PostedAt,
	}
}
