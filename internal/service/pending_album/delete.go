package pending_album

import (
	"context"
)

func (s *service) Delete(ctx context.Context, id int) error {
	err := s.pendingAlbumRepo.DeleteByAlbumId(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
