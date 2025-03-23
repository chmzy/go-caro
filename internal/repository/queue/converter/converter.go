package converter

import (
	modelrepo "go-caro/internal/repository/queue/model"
	modelserv "go-caro/internal/service/queue/model"
)

func ToPostQueueFromRepo(post *modelrepo.PostQueue) *modelserv.PostQueue {
	return &modelserv.PostQueue{
		ID:      post.ID,
		Author:  post.Author,
		AlbumID: post.AlbumID,
		MsgLink: modelserv.Link(post.MsgLink),
	}
}
