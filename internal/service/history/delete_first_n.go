package history

import "context"

func (s *service) DeleteFirstN(ctx context.Context, n uint64) error {
	//TODO: add transactions
	err := s.historyRepo.DeleteFirstN(ctx, n)
	if err != nil {
		return err
	}

	return nil
}
