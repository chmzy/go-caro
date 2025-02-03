package api

import (
	"context"
	"fmt"
	c "go-caro/internal/service/queue/converter"
	mw "go-caro/pkg/tg/middleware"
	m "go-caro/pkg/tg/model"
	"log"
)

func (a *API) OnMedia(ctx m.Context) error {
	adminUsers := ctx.Get("admins").([]string)

	return mw.FromAdmin(adminUsers, a.onMediaAdmin, a.onMediaUser)(ctx)
}

func (a *API) onMediaAdmin(ctx m.Context) error {
	chanId := ctx.Get("chan_id").(int64)

	deletePost := func(ctx m.Context) error {
		if err := a.historyService.DeleteByID(context.Background(), uint64(chanId)); err != nil {
			return err
		}
		ctx.Bot().Delete(m.Post{
			MessageID: fmt.Sprintf("%d", ctx.Message().OriginalMessageID),
			ChatID:    ctx.Message().OriginalChat.ID,
		})

		log.Println("Deleted post from channel")

		return nil
	}

	saveMediaMsg := func(ctx m.Context) error {
		id, err := a.queueService.Put(context.Background(), c.ToQueueFromAPI(ctx.Message()))
		if err != nil {
			return err
		}

		if err := ctx.Send(fmt.Sprintf("Thanks for media, admin! Saved with id %d", id)); err != nil {
			return err
		}

		return nil
	}

	return mw.ForwardedFromChannel(chanId, deletePost, saveMediaMsg)(ctx)
}

func (a *API) onMediaUser(ctx m.Context) error {
	if err := ctx.Send("Thx for media!"); err != nil {
		return err
	}

	return nil
}
