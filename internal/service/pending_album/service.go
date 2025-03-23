package pending_album

import (
	repo "go-caro/internal/repository"
	serv "go-caro/internal/service"
)

type service struct {
	pendingAlbumRepo repo.PendingAlbumRepository
}

func NewService(repo repo.PendingAlbumRepository) serv.PendingAlbumService {
	return &service{
		pendingAlbumRepo: repo,
	}
}
