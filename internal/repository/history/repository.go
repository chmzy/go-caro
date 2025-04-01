package history

import (
	"context"
	"fmt"
	"go-caro/internal/repository"
	"go-caro/internal/repository/history/converter"
	modelrepo "go-caro/internal/repository/history/model"
	modelserv "go-caro/internal/service/history/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName      = "history"
	idColumn       = "id"
	albumIdColumn  = "album_id"
	chatIdColumn   = "chat_id"
	msgIdColumn    = "msg_id"
	postedAtColumn = "posted_at"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.HistoryRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, post *modelserv.PostHistory) (int, error) {
	var id int
	sql := fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s) VALUES ($1,$2,$3,$4) RETURNING %s", tableName, albumIdColumn, chatIdColumn, msgIdColumn, postedAtColumn, idColumn)
	err := r.db.QueryRow(ctx, sql, post.AlbumID, post.ChatID, post.MsgID, post.PostedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo: history: Create: %w", err)
	}

	// When write into database, we pass external type struct from model package
	// No need to convert
	return id, nil
}

func (r *repo) GetLast(ctx context.Context) (*modelserv.PostHistory, error) {
	var post modelrepo.PostHistory
	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY %s DESC LIMIT 1", tableName, postedAtColumn)
	err := r.db.QueryRow(ctx, sql).Scan(&post.ID, &post.AlbumID, &post.ChatID, &post.MsgID, &post.PostedAt)
	if err != nil {
		return nil, fmt.Errorf("repo: history: GetLast: %w", err)
	}
	// When fetch data from table
	// We parse data into table struct
	// Then convert this struct into external type from model package
	return converter.ToHistoryFromRepo(&post), nil
}

func (r *repo) DeleteByID(ctx context.Context, id int) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, idColumn)
	_, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("repo: history: DeleteByID: %w", err)
	}

	return nil
}

func (r *repo) DeleteByAlbumID(ctx context.Context, id string) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, albumIdColumn)
	_, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("repo: history: DeleteByAlbumID: %w", err)
	}

	return nil
}

func (r *repo) DeleteKeepLastN(ctx context.Context, n uint64) error {
	query := `
  	WITH ranked_rows AS (
  	SELECT 
  	  *,
  	  ROW_NUMBER() OVER (ORDER BY id DESC) AS row_num
  	FROM %s
  	)
  	DELETE FROM %s
  	WHERE id IN (
  	  SELECT id 
  	  FROM ranked_rows 
  	  WHERE row_num > 10
  	);
  `
	sql := fmt.Sprintf(query, tableName, tableName)
	_, err := r.db.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("repo: history: DeleteKeepLastN: %w", err)
	}

	return nil
}
