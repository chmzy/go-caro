package pending_album

import (
	"context"
	"fmt"
	"go-caro/internal/repository"
	"log"

	c "go-caro/internal/repository/pending_album/converter"
	modelrepo "go-caro/internal/repository/pending_album/model"
	modelserv "go-caro/internal/service/pending_album/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	tableName     = "pending_album"
	idColumn      = "id"
	albumIdColumn = "album_id"
	authorColumn  = "author"
	chatIdColumn  = "chat_id"
	msgIdColumn   = "msg_id"
)

type repo struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) repository.PendingAlbumRepository {
	return &repo{
		db: db,
	}
}

func (r *repo) Put(ctx context.Context, album *modelserv.AlbumPost) error {
	sql := fmt.Sprintf("INSERT INTO %s (%s,%s,%s,%s) VALUES ($1, $2, $3, $4)", tableName, albumIdColumn, authorColumn, chatIdColumn, msgIdColumn)
	log.Println(album.AlbumID, album.Author, album.MsgLink.ChatID, album.MsgLink.MsgID)
	_, err := r.db.Exec(ctx, sql, album.AlbumID, album.Author, album.MsgLink.ChatID, album.MsgLink.MsgID)
	if err != nil {
		return fmt.Errorf("repo: pending_album: put: %w", err)
	}

	return nil
}

func (r *repo) Next(ctx context.Context) ([]modelserv.AlbumPost, error) {
	sql := fmt.Sprintf("SELECT * FROM %s WHERE %s = $1", tableName, albumIdColumn)
	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, fmt.Errorf("repo: pending_album: next: %w", err)
	}
	defer rows.Close()

	var posts []modelrepo.AlbumPost
	posts, err = pgx.CollectRows(rows, pgx.RowToStructByName[modelrepo.AlbumPost])

	return c.ToPendingAlbumFromRepo(posts), nil
}

func (r *repo) DeleteByAlbumId(ctx context.Context, id int) error {
	sql := fmt.Sprintf("DELETE FROM %s WHERE %s = $1", tableName, albumIdColumn)
	if _, err := r.db.Exec(ctx, sql, id); err != nil {
		return fmt.Errorf("repo: pending_album: delete: %w", err)
	}

	return nil
}
