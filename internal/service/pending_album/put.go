package pending_album

import (
	"context"
	modelserv "go-caro/internal/service/pending_album/model"
)

func (s *service) Put(ctx context.Context, post *modelserv.AlbumPost) error {
	err := s.pendingAlbumRepo.Put(ctx, post)
	if err != nil {
		return err
	}
	return nil
}
