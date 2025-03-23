package api

import (
	"go-caro/internal/service"
)

type API struct {
	historyService      service.HistoryService
	queueService        service.QueueService
	pendingAlbumService service.PendingAlbumService
}

func NewAPI(hs service.HistoryService, qs service.QueueService, pas service.PendingAlbumService) *API {
	return &API{
		historyService:      hs,
		queueService:        qs,
		pendingAlbumService: pas,
	}
}
