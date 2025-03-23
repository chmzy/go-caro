package api

import (
	"context"
	"strings"
	"sync"
	"time"

	"fmt"
	c "go-caro/internal/service/queue/converter"
	mw "go-caro/pkg/tg/middleware"
	m "go-caro/pkg/tg/model"
	"log"

	"gopkg.in/telebot.v4"
)

var (
	pendingAlbums map[string]chan *m.Message = make(map[string]chan *m.Message, 2)
	mu            sync.Mutex
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
	// Define stored message
	msg := telebot.StoredMessage{
		ChatID:    ctx.Chat().ID,
		MessageID: fmt.Sprintf("%d", ctx.Message().ID),
	}

	if ctx.Message().AlbumID != "" {
		aId := ctx.Message().AlbumID
		mu.Lock()
		defer mu.Unlock()

		// Check if an existing goroutine is collecting messages
		ch, exists := pendingAlbums[aId]
		if !exists {
			// Create a new channel for this AlbumID
			dataChan := make(chan *m.Message, 10)
			pendingAlbums[aId] = dataChan

			// Spawn a new goroutine to handle messages
			go processAlbum(ctx, aId, dataChan)
			ch = dataChan
		}

		ch <- ctx.Message()
		return nil
	}
	var (
		inlineKeys = &telebot.ReplyMarkup{}
		btnApply   = inlineKeys.Data("✅ Apply", "apply_action")
		btnReject  = inlineKeys.Data("❌ Reject", "reject_action")
	)
	inlineKeys.Inline(telebot.Row{btnApply, btnReject})

	_, err := ctx.Bot().Copy(&telebot.Chat{ID: -1002504066662}, msg, &telebot.SendOptions{
		ReplyMarkup: inlineKeys,
	})

	if err != nil {
		return err
	}

	return nil
}

func processAlbum(bot m.Context, albumID string, dataChan chan *m.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer delete(pendingAlbums, albumID) // Cleanup after processing
	defer close(dataChan)
	defer cancel()

	var messages []*m.Message = make([]*m.Message, 0, 2)
	for {
		select {
		case <-ctx.Done():
			var g []telebot.Inputtable
			for _, m := range messages {
				switch strings.ToLower(m.Media().MediaType()) {
				case "photo":
					g = append(g, m.Photo)
				case "video":
					g = append(g, m.Video)
				case "animation", "gif":
					g = append(g, m.Animation)
				default:
					log.Println("onMediaUser: processAlbum: unsupported album item", m.Media().MediaType())
					return
				}
			}

			// 1. Send the album first
			msgs, err := bot.Bot().SendAlbum(&telebot.Chat{ID: -1002504066662}, g)
			if err != nil {
				log.Printf("Failed to send album: %s", err.Error())
				return
			}

			// 2. Send a separate message with the inline keyboard
			inlineKeys := &telebot.ReplyMarkup{}
			btnApply := inlineKeys.Data("✅ Apply", "apply_action")
			btnReject := inlineKeys.Data("❌ Reject", "reject_action")
			inlineKeys.Inline(telebot.Row{btnApply, btnReject})

			_, err = bot.Bot().Send(
				&telebot.Chat{ID: -1002504066662},
				fmt.Sprintf("%d", len(msgs)),
				&telebot.SendOptions{ReplyMarkup: inlineKeys, ReplyTo: &msgs[0]},
			)
			if err != nil {
				log.Printf("Failed to send keyboard: %s", err.Error())
				return
			}
			return
		case msg := <-dataChan:
			messages = append(messages, msg)
		}
	}
}
