package history

import (
	"context"
)

func (s *service) DeleteByAlbumID(ctx context.Context, id string) error {
	//TODO: add transactions
	err := s.historyRepo.DeleteByAlbumID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
