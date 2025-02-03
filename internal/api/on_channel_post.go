package api

import (
	"context"
	c "go-caro/internal/service/history/converter"
	m "go-caro/pkg/tg/model"
	"log"
)

func (a *API) OnChannelPost(ctx m.Context) error {
	if msg := ctx.Message().Media(); msg != nil {
		switch msg.MediaType() {
		case "photo", "video", "gif":
			id, err := a.historyService.Create(context.Background(), c.ToHistoryFromAPI(ctx.Message()))
			if err != nil {
				return err
			}
			log.Printf("Add post with id %d", id)
		default:
			break
		}
	}

	return nil

}
