package api

import (
	"context"
	"fmt"
	"go-caro/pkg/math"
	m "go-caro/pkg/tg/model"
)

const (
	tgMsgLink = "https://t.me/c/%d/%s"
)

func (a *API) Queue(ctx m.Context) error {
	posts, err := a.queueService.Next(context.Background())
	if err != nil {
		return err
	}

	chatID := math.AbsInt64(posts[0].ChatID + 1000000000000)
	msg := fmt.Sprintf(tgMsgLink, chatID, posts[0].MsgID)

	if err = ctx.Reply(msg); err != nil {
		return fmt.Errorf("queue: reply: %w\n", err)
	}

	return nil
}
