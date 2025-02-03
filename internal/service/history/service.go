package history

import (
	repo "go-caro/internal/repository"
	serv "go-caro/internal/service"
)

type service struct {
	historyRepo repo.HistoryRepository
}

func NewService(repo repo.HistoryRepository) serv.HistoryService {
	return &service{
		historyRepo: repo,
	}
}
