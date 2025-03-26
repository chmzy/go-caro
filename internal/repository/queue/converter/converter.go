package converter

import (
	modelrepo "go-caro/internal/repository/queue/model"
	modelserv "go-caro/internal/service/queue/model"
)

func ToPostQueueFromRepo(posts []modelrepo.PostQueue) []modelserv.PostQueue {
	var p []modelserv.PostQueue = make([]modelserv.PostQueue, 0, len(posts))
	for _, post := range posts {
		p = append(p, modelserv.PostQueue{
			ID:      post.ID,
			Author:  post.Author,
			AlbumID: post.AlbumID,
			ChatID:  post.ChatID,
			MsgID:   post.MsgID,
		})
	}
	return p
}
