package pending_album

import (
	"context"
	modelserv "go-caro/internal/service/pending_album/model"
)

func (s *service) Next(ctx context.Context) ([]modelserv.AlbumPost, error) {
	posts, err := s.pendingAlbumRepo.Next(ctx)
	if err != nil {
		return nil, err
	}
	return posts, nil
}
