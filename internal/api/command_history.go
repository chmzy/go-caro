package api

import (
	"context"
	"fmt"
	"go-caro/pkg/math"
	m "go-caro/pkg/tg/model"
)

func (a *API) History(ctx m.Context) error {

	post, err := a.historyService.GetLast(context.Background())
	if err != nil {
		return err
	}
	chatID := math.AbsInt64(post.ChatID + 1000000000000)
	msg := fmt.Sprintf(tgMsgLink, chatID, post.MsgID)

	if err := ctx.Reply(msg); err != nil {
		return err
	}

	return nil
}
