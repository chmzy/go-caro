package flusher

import (
	"context"
	"go-caro/internal/repository"
	"log"
	"time"
)

const (
	nOfPosts = 5
)

type flusher struct {
	hisotryRepo repository.HistoryRepository
}

func NewFlusher(historyRepo repository.HistoryRepository) *flusher {
	return &flusher{
		hisotryRepo: historyRepo,
	}
}

func (f *flusher) Start(ctx context.Context, period int64) {
	p := time.Duration(period) * time.Second
	ticker := time.NewTicker(p)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				//noop
			}
			if err := f.hisotryRepo.DeleteKeepLastN(context.Background(), nOfPosts); err != nil {
				log.Println("flusher: historyRepo.DeleteKeepLastN: %w\n", err)
			}

			log.Println("Flushed history successfuly.")
			continue
		}
	}()
}
