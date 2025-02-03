package queue

import (
	"context"
	modelserv "go-caro/internal/service/queue/model"
)

func (s *service) Put(ctx context.Context, post *modelserv.PostQueue) (int, error) {
	id, err := s.queueRepo.Put(ctx, post)
	if err != nil {
		return 0, err
	}
	return id, nil
}
