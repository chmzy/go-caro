package api

import (
	"context"
	"fmt"
	"go-caro/internal/events"
	"go-caro/internal/service/queue/model"
	m "go-caro/pkg/tg/model"
	"log"
	"strconv"

	"gopkg.in/telebot.v4"
)

func (a *API) ApplyButtonEvent(ctx m.Context) error {
	// Check if message is media album
	if ctx.Message().ReplyTo != nil && ctx.Message().ReplyTo.AlbumID != "" {
		return a.applyAlbum(ctx)
	}

	return a.applySingle(ctx)
}

func (a *API) applySingle(ctx m.Context) error {
	_, err := a.queueService.Put(context.Background(), &model.PostQueue{
		Author:  ctx.Message().OriginalSenderName,
		AlbumID: "",
		ChatID:  ctx.Message().Chat.ID,
		MsgID:   fmt.Sprintf("%d", ctx.Message().ID),
	})
	if err != nil {
		log.Println(err)
	}

	return editMarkupKeyboard(ctx)
}

func (a *API) applyAlbum(ctx m.Context) error {
	originalMessage := ctx.Message().ReplyTo
	albumLen, _ := strconv.Atoi(ctx.Message().Text)
	for i := range albumLen {
		_, err := a.queueService.Put(context.Background(), &model.PostQueue{
			Author:  originalMessage.OriginalSenderName,
			AlbumID: originalMessage.AlbumID,
			ChatID:  originalMessage.Chat.ID,
			MsgID:   fmt.Sprintf("%d", originalMessage.ID+i),
		})

		if err != nil {
			log.Println(err)
		}

	}

	return editMarkupKeyboard(ctx)
}

func editMarkupKeyboard(ctx m.Context) error {
	inlineKeys := &telebot.ReplyMarkup{}
	btnApproved := inlineKeys.Data("⏳ Will be posted soon...", "noop")
	btnDelete := inlineKeys.Data("🚫 Delete from queue", events.DeleteButton)
	inlineKeys.Inline(telebot.Row{btnApproved}, telebot.Row{btnDelete})

	err := ctx.Edit(inlineKeys)
	if err != nil {
		return err
	}

	return nil

}
