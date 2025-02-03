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

func (r *repo) Create(ctx context.Context, post *modelserv.PostHistory) (uint64, error) {
	var id uint64
	sql := fmt.Sprintf("INSERT INTO %s (%s,%s) VALUES ($1,$2) RETURNING %s", tableName, idColumn, postedAtColumn, idColumn)
	err := r.db.QueryRow(ctx, sql, post.ID, post.PostedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	// When write into database, we pass external type struct from model package
	// No need to convert
	return id, nil
}

func (r *repo) GetLast(ctx context.Context) (*modelserv.PostHistory, error) {
	var post modelrepo.PostHistory
	sql := fmt.Sprintf("SELECT (%s,%s) FROM %s ORDER BY %s DESC LIMIT 1", idColumn, postedAtColumn, tableName, postedAtColumn)
	err := r.db.QueryRow(ctx, sql).Scan(&post)
	if err != nil {
		return nil, err
	}
	// When fetch data from table
	// We parse data into table struct
	// Then convert this struct into external type from model package
	return converter.ToHistoryFromRepo(&post), nil
}

func (r *repo) DeleteByID(ctx context.Context, id uint64) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = %d", tableName, idColumn, id)
	_, err := r.db.Exec(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) DeleteFirstN(ctx context.Context, n uint64) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s IN (SELECT %s FROM %s ORDER BY %s LIMIT %d);", tableName, idColumn, idColumn, tableName, idColumn, n)
	_, err := r.db.Exec(ctx, sql)
	if err != nil {
		return err
	}
	return nil
}
