package queue

import (
	"context"
)

func (s *service) DeleteByMsgID(ctx context.Context, id string) error {
	err := s.queueRepo.DeleteByMsgID(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
