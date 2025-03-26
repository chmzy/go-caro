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
	err := r.db.QueryRow(ctx, sql, post.Author, post.AlbumID, post.ChatID, post.MsgID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo: queue: put: %w", err)
	}

	return id, nil
}

func (r *repo) Next(ctx context.Context) ([]modelserv.PostQueue, error) {
	var posts []modelrepo.PostQueue
	query := `
	WITH first_row AS (
	    SELECT * 
	    FROM %s
	    ORDER BY %s
	    LIMIT 1
	), grouped_rows AS (
	    SELECT * 
	    FROM %s
	    WHERE %s = (SELECT %s FROM first_row) AND %s <> ''
	    UNION ALL
	    SELECT * FROM first_row 
	    WHERE %s = ''
	)
	SELECT * FROM grouped_rows
	ORDER BY %s;
	`
	sql := fmt.Sprintf(query, tableName, idColumn, tableName, albumIdColumn, albumIdColumn, albumIdColumn, albumIdColumn, idColumn)
	rows, err := r.db.Query(ctx, sql)
	defer rows.Close()

	posts, err = pgx.CollectRows(rows, pgx.RowToStructByNameLax[modelrepo.PostQueue])
	if err != nil {
		return nil, fmt.Errorf("repo: queue: next: %w", err)
	}

	return c.ToPostQueueFromRepo(posts), nil
}

func (r *repo) DeleteByMsgID(ctx context.Context, id string) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s=$1", tableName, msgIdColumn)
	if _, err := r.db.Exec(ctx, sql, id); err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteByAlbumID(ctx context.Context, id string) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s=$1", tableName, albumIdColumn)
	if _, err := r.db.Exec(ctx, sql, id); err != nil {
		return err
	}

	return nil
}
