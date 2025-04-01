package history

import "context"

func (s *service) DeleteKeepLastN(ctx context.Context, n uint64) error {
	//TODO: add transactions
	err := s.historyRepo.DeleteKeepLastN(ctx, n)
	if err != nil {
		return err
	}

	return nil
}
