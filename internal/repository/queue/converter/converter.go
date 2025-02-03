package converter

import (
	modelrepo "go-caro/internal/repository/queue/model"
	modelserv "go-caro/internal/service/queue/model"
)

func ToQueuePostFromRepo(post *modelrepo.PostQueue) *modelserv.PostQueue {
	return &modelserv.PostQueue{
		ID:      post.ID,
		MsgLink: modelserv.Link(post.MsgLink),
	}
}
