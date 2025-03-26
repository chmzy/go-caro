package history

import (
	"context"
	modelserv "go-caro/internal/service/history/model"
)

func (s *service) Create(ctx context.Context, post *modelserv.PostHistory) (int, error) {
	//TODO: add transactions
	id, err := s.historyRepo.Create(ctx, post)
	if err != nil {
		return 0, err
	}

	return id, nil
}
