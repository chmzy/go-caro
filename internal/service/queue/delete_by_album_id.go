package queue

import (
	"context"
)

func (s *service) DeleteByAlbumID(ctx context.Context, id string) error {
	err := s.queueRepo.DeleteByAlbumID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
