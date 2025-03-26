package api

import (
	"fmt"
	m "go-caro/pkg/tg/model"
	"strconv"

	"gopkg.in/telebot.v4"
)

func (a *API) RejectButtonEvent(ctx m.Context) error {
	// Check if message is media album
	if ctx.Message().ReplyTo != nil && ctx.Message().ReplyTo.AlbumID != "" {
		return a.rejectAlbum(ctx)
	}

	return a.rejectSingle(ctx)
}

func (a *API) rejectAlbum(ctx m.Context) error {
	originalMessage := ctx.Message().ReplyTo
	albumLen, _ := strconv.Atoi(ctx.Message().Text)
	var albumMsgs []telebot.Editable
	for i := range albumLen {
		albumMsgs = append(albumMsgs, telebot.StoredMessage{ChatID: originalMessage.Chat.ID, MessageID: fmt.Sprintf("%d", originalMessage.ID+i)})
	}

	if err := ctx.Bot().DeleteMany(albumMsgs); err != nil {
		return err
	}

	return a.rejectSingle(ctx)
}

func (a *API) rejectSingle(ctx m.Context) error {
	// Delete message with inline keyboard
	if err := ctx.Delete(); err != nil {
		return err
	}
	return nil
}
