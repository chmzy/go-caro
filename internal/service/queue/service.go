package queue

import (
	repo "go-caro/internal/repository"
	serv "go-caro/internal/service"
)

type service struct {
	queueRepo repo.QueueRepository
}

func NewService(repo repo.QueueRepository) serv.QueueService {
	return &service{
		queueRepo: repo,
	}
}
