package history

import (
	"context"
	"go-caro/internal/service/history/model"
)

func (s *service) GetLast(ctx context.Context) (*model.PostHistory, error) {
	//TODO: add transactions
	post, err := s.historyRepo.GetLast(ctx)
	if err != nil {
		return nil, err
	}

	return post, nil
}
