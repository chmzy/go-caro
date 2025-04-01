package api

import (
	"go-caro/internal/service"
)

type API struct {
	historyService service.HistoryService
	queueService   service.QueueService
}

func NewAPI(hs service.HistoryService, qs service.QueueService) *API {
	return &API{
		historyService: hs,
		queueService:   qs,
	}
}
