package queue

import (
	"context"
	modelserv "go-caro/internal/service/queue/model"
)

func (s *service) Next(ctx context.Context) ([]modelserv.PostQueue, error) {
	posts, err := s.queueRepo.Next(ctx)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
