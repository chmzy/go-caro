package queue

import (
	"context"
	"fmt"
	"go-caro/internal/repository"
	c "go-caro/internal/repository/queue/converter"

	modelrepo "go-caro/internal/repository/queue/model"
	modelserv "go-caro/internal/service/queue/model"

	"github.com/jackc/pgx/v5/pgxpool"
	// modelrepo "go-caro/internal/repository/queue/model"
)

const (
	tableName    = "queue"
	idColumn     = "id"
	chatIdColumn = "chat_id"
	msgIdColumn  = "msg_id"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.QueueRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Put(ctx context.Context, post *modelserv.PostQueue) (int, error) {
	var id int
	sql := fmt.Sprintf("INSERT INTO %s (%s,%s) VALUES ($1,$2) RETURNING %s", tableName, chatIdColumn, msgIdColumn, idColumn)
	err := r.db.QueryRow(ctx, sql, post.MsgLink.ChatID, post.MsgLink.MsgID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo: queue: put: %w", err)
	}

	return id, nil
}

func (r *repo) Next(ctx context.Context) (*modelserv.PostQueue, error) {
	var post modelrepo.PostQueue
	sql := fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
	err := r.db.QueryRow(ctx, sql).Scan(&post.ID, &post.MsgLink.ChatID, &post.MsgLink.MsgID)
	if err != nil {
		return nil, fmt.Errorf("repo: queue: put: %w", err)
	}

	return c.ToQueuePostFromRepo(&post), nil
}

func (r *repo) Delete(ctx context.Context, id int) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s=$1", tableName, idColumn)
	if _, err := r.db.Exec(ctx, sql, id); err != nil {
		return err
	}

	return nil
}
