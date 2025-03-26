package api

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"fmt"
	"go-caro/internal/events"
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
		if err := a.historyService.DeleteByID(context.Background(), ctx.Message().OriginalMessageID); err != nil {
			return err
		}

		err := ctx.Bot().Delete(m.Post{
			MessageID: fmt.Sprintf("%d", ctx.Message().OriginalMessageID),
			ChatID:    ctx.Message().OriginalChat.ID,
		})
		if err != nil {
			return err
		}

		if err := ctx.Delete(); err != nil {
			return err
		}

		log.Println("Deleted post from channel")

		return nil
	}

	saveMediaMsg := func(ctx m.Context) error {
		log.Println("New message! ID: ", ctx.Message().ID)
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
		mu.Lock()
		defer mu.Unlock()
		aId := ctx.Message().AlbumID
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
		btnApply   = inlineKeys.Data("✅ Apply", events.ApplyButton)
		btnReject  = inlineKeys.Data("❌ Reject", events.RejectButton)
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer delete(pendingAlbums, albumID) // Cleanup after processing
	defer close(dataChan)
	defer cancel()

	var messages []*m.Message = make([]*m.Message, 0, 2)
	for {
		select {
		case <-ctx.Done():
			sort.Slice(messages, func(i, j int) bool {
				return messages[i].ID <= messages[j].ID
			})

			var g telebot.Album
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
			g.SetCaption("[Caro est infirma](https://t.me/caroinfirma) ❤️ [Suggest a post](https://t.me/Caro_est_infirma_bot)")
			msgs, err := bot.Bot().SendAlbum(&telebot.Chat{ID: -1002504066662}, g, &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
			if err != nil {
				log.Printf("Failed to send album: %s", err.Error())
				return
			}

			// 2. Send a separate message with the inline keyboard
			inlineKeys := &telebot.ReplyMarkup{}
			btnApply := inlineKeys.Data("✅ Apply", events.ApplyButton)
			btnReject := inlineKeys.Data("❌ Reject", events.RejectButton)
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
