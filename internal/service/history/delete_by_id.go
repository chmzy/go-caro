package history

import (
	"context"
)

func (s *service) DeleteByID(ctx context.Context, id int) error {
	//TODO: add transactions
	err := s.historyRepo.DeleteByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
