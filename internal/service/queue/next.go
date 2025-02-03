package queue

import (
	"context"
	modelserv "go-caro/internal/service/queue/model"
)

func (s *service) Next(ctx context.Context) (*modelserv.PostQueue, error) {
	post, err := s.queueRepo.Next(ctx)
	if err != nil {
		return nil, err
	}

	return post, nil
}
