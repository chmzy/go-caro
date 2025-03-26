package api

import (
	"context"
	"fmt"
	m "go-caro/pkg/tg/model"
	"strconv"

	"gopkg.in/telebot.v4"
)

func (a *API) DeleteButtonEvent(ctx m.Context) error {
	// Check if message is media album
	if ctx.Message().ReplyTo != nil && ctx.Message().ReplyTo.AlbumID != "" {
		return a.deleteAlbum(ctx)
	}

	return a.deleteSingle(ctx)
}

func (a *API) deleteAlbum(ctx m.Context) error {
	originalMessage := ctx.Message().ReplyTo

	if err := a.queueService.DeleteByAlbumID(context.Background(), originalMessage.AlbumID); err != nil {
		return err
	}

	albumLen, _ := strconv.Atoi(ctx.Message().Text)
	var albumMsgs []telebot.Editable
	for i := range albumLen {
		albumMsgs = append(albumMsgs, telebot.StoredMessage{ChatID: originalMessage.Chat.ID, MessageID: fmt.Sprintf("%d", originalMessage.ID+i)})
	}

	if err := ctx.Bot().DeleteMany(albumMsgs); err != nil {
		return err
	}

	if err := ctx.Delete(); err != nil {
		return err
	}

	return nil
}

func (a *API) deleteSingle(ctx m.Context) error {
	if err := a.queueService.DeleteByMsgID(context.Background(), fmt.Sprintf("%d", ctx.Message().ID)); err != nil {
		return err
	}

	// Delete message with inline keyboard
	if err := ctx.Delete(); err != nil {
		return err
	}
	return nil
}
