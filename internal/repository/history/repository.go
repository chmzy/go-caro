package history

import (
	"context"
	"fmt"
	"go-caro/internal/repository"
	"go-caro/internal/repository/history/converter"
	modelrepo "go-caro/internal/repository/history/model"
	modelserv "go-caro/internal/service/history/model"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName      = "history"
	idColumn       = "id"
	albumIdColumn  = "album_id"
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
	sql := fmt.Sprintf("INSERT INTO %s VALUES ($1,$2,$3) RETURNING %s", tableName, idColumn)
	err := r.db.QueryRow(ctx, sql, post.ID, post.AlbumID, post.PostedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("repo: history: create: %w", err)
	}

	// When write into database, we pass external type struct from model package
	// No need to convert
	return id, nil
}

func (r *repo) GetLast(ctx context.Context) (*modelserv.PostHistory, error) {
	var post modelrepo.PostHistory
	sql := fmt.Sprintf("SELECT * FROM %s ORDER BY %s DESC LIMIT 1", tableName, postedAtColumn)
	err := r.db.QueryRow(ctx, sql).Scan(&post.ID, &post.AlbumID, &post.PostedAt)
	if err != nil {
		return nil, fmt.Errorf("repo: history: getLast: %w", err)
	}
	// When fetch data from table
	// We parse data into table struct
	// Then convert this struct into external type from model package
	return converter.ToHistoryFromRepo(&post), nil
}

func (r *repo) DeleteByID(ctx context.Context, id int) error {
	log.Println(id)
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, idColumn)
	_, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("repo: history: deleteByID: %w", err)
	}

	return nil
}

func (r *repo) DeleteByAlbumID(ctx context.Context, id string) error {
	log.Println(id)
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, albumIdColumn)
	_, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("repo: history: deleteByAlbumID: %w", err)
	}

	return nil
}

func (r *repo) DeleteFirstN(ctx context.Context, n uint64) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s IN (SELECT %s FROM %s ORDER BY %s LIMIT %d);", tableName, idColumn, idColumn, tableName, idColumn, n)
	_, err := r.db.Exec(ctx, sql)
	if err != nil {
		return fmt.Errorf("repo: history: deleteFirstN: %w", err)
	}
	return nil
}
