package repository

import (
	"context"
	hmodelserv "go-caro/internal/service/history/model"
	pamodelserv "go-caro/internal/service/pending_album/model"
	qmodelserv "go-caro/internal/service/queue/model"
)

type (
	HistoryRepository interface {
		// Add post
		Create(ctx context.Context, post *hmodelserv.PostHistory) (int, error)
		// Returns last post
		GetLast(ctx context.Context) (*hmodelserv.PostHistory, error)
		// Delete post by id
		DeleteByID(ctx context.Context, id int) error
		// Delete album by id
		DeleteByAlbumID(ctx context.Context, id string) error
		// Delets N posts from the table's begining
		DeleteFirstN(ctx context.Context, n uint64) error
	}

	QueueRepository interface {
		// Put post
		Put(ctx context.Context, post *qmodelserv.PostQueue) (int, error)
		// Get next post or album
		Next(ctx context.Context) ([]qmodelserv.PostQueue, error)
		// Delete post by id
		DeleteByMsgID(ctx context.Context, id string) error
		// Delete album by id
		DeleteByAlbumID(ctx context.Context, id string) error
	}

	PendingAlbumRepository interface {
		// Put album post
		Put(ctx context.Context, post *pamodelserv.AlbumPost) error
		// Get next album posts
		Next(ctx context.Context) ([]pamodelserv.AlbumPost, error)
		// Delete album posts
		DeleteByAlbumId(ctx context.Context, id int) error
	}
)
