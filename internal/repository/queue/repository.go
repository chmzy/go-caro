package queue

import (
	"context"
	"fmt"
	"go-caro/internal/repository"
	c "go-caro/internal/repository/queue/converter"

	modelrepo "go-caro/internal/repository/queue/model"
	modelserv "go-caro/internal/service/queue/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName     = "queue"
	idColumn      = "id"
	authorColumn  = "author"
	albumIdColumn = "album_id"
	chatIdColumn  = "chat_id"
	msgIdColumn   = "msg_id"
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
	sql := fmt.Sprintf("INSERT INTO %s (%s, %s, %s,%s) VALUES ($1, $2, $3, $4) RETURNING %s", tableName, authorColumn, albumIdColumn, chatIdColumn, msgIdColumn, idColumn)
	err := r.db.QueryRow(ctx, sql, post.Author, post.AlbumID, post.MsgLink.ChatID, post.MsgLink.MsgID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo: queue: put: %w", err)
	}

	return id, nil
}

func (r *repo) Next(ctx context.Context) (*modelserv.PostQueue, error) {
	var post modelrepo.PostQueue
	sql := fmt.Sprintf("SELECT * FROM %s LIMIT 1", tableName)
	rows, err := r.db.Query(ctx, sql)
	defer rows.Close()

	post, err = pgx.CollectOneRow(rows, pgx.RowToStructByName[modelrepo.PostQueue])
	if err != nil {
		return nil, fmt.Errorf("repo: queue: put: %w", err)
	}

	return c.ToPostQueueFromRepo(&post), nil
}

func (r *repo) DeleteById(ctx context.Context, id int) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s=$1", tableName, idColumn)
	if _, err := r.db.Exec(ctx, sql, id); err != nil {
		return err
	}

	return nil
}
func (r *repo) DeleteByAlbumId(ctx context.Context, id int) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s=$1", tableName, albumIdColumn)
	if _, err := r.db.Exec(ctx, sql, id); err != nil {
		return err
	}

	return nil
}
