package api

import (
	"context"
	"slices"
	"sort"
	"strings"
	"sync"
	"time"

	"fmt"
	"go-caro/internal/events"
	"go-caro/internal/service/queue/converter"
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

	sendPostToQueue := func(ctx m.Context) error {
		keyboard := &telebot.ReplyMarkup{}
		btnApproved := keyboard.Data("⏳ Will be posted soon...", "noop")
		btnDelete := keyboard.Data("🚫 Delete from queue", events.DeleteButton)
		keyboard.Inline(telebot.Row{btnApproved}, telebot.Row{btnDelete})

		afterFn := func(msg *m.Message) error {
			id, err := a.queueService.Put(context.Background(), converter.ToQueueFromAPI(msg))
			if err != nil {
				return err
			}
			if err := ctx.Send(fmt.Sprintf("Thanks for media, admin! Saved with id %d", id)); err != nil {
				return err
			}

			return nil
		}

		if ctx.Message().AlbumID == "" {
			msg, err := sendSingle(ctx, keyboard)
			if err != nil {
				return err
			}
			if err := afterFn(msg); err != nil {
				return err
			}

			return nil
		}

		sendAlbum(ctx, keyboard, afterFn)

		return nil
	}

	return mw.ForwardedFromChannel(chanId, deletePost, sendPostToQueue)(ctx)
}

func (a *API) onMediaUser(ctx m.Context) error {
	keyboard := &telebot.ReplyMarkup{}
	btnApply := keyboard.Data("✅ Apply", events.ApplyButton)
	btnReject := keyboard.Data("❌ Reject", events.RejectButton)
	keyboard.Inline(telebot.Row{btnApply, btnReject})

	if ctx.Message().AlbumID == "" {
		if _, err := sendSingle(ctx, keyboard); err != nil {
			return err
		}
		return nil
	}
	sendAlbum(ctx, keyboard, nil)

	return nil
}

func sendSingle(ctx m.Context, keyboard *telebot.ReplyMarkup) (*m.Message, error) {
	suggestChanId := ctx.Get("suggest_id").(int64)
	admins := ctx.Get("admins").([]string)
	media, err := createMediaItem(ctx.Message())
	if err != nil {
		return nil, err
	}
	a := telebot.Album{media}

	senderName := ctx.Sender().FirstName
	if slices.Contains(admins, ctx.Sender().Username) {
		a.SetCaption("[Caro est infirma](https://t.me/caroinfirma) ❤️ [Suggest a post](https://t.me/Caro_est_infirma_bot)")
	} else {
		caption := fmt.Sprintf("From %s \n\n [Caro est infirma](https://t.me/caroinfirma) ❤️ [Suggest a post](https://t.me/Caro_est_infirma_bot)", senderName)
		a.SetCaption(caption)
	}

	msg, err := ctx.Bot().Send(&telebot.Chat{ID: suggestChanId}, a[0], &telebot.SendOptions{
		ParseMode:   telebot.ModeMarkdown,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func sendAlbum(ctx m.Context, keyboard *telebot.ReplyMarkup, afterFn func(m *m.Message) error) {
	// Handle album case
	mu.Lock()
	defer mu.Unlock()

	albumID := ctx.Message().AlbumID
	ch, exists := pendingAlbums[albumID]
	if !exists {
		ch = make(chan *m.Message, 10)
		pendingAlbums[albumID] = ch
		go func() {
			msgs, err := processAlbum(ctx, albumID, ch, keyboard)
			if err != nil {
				log.Println(err)
				return
			}

			if afterFn != nil {
				for _, m := range msgs {
					if err := afterFn(&m); err != nil {
						log.Println(err)
					}
				}
			}

			return
		}()
	}

	ch <- ctx.Message()

}

// Helper function to create a media item from a message
func createMediaItem(msg *m.Message) (telebot.Inputtable, error) {
	switch strings.ToLower(msg.Media().MediaType()) {
	case "photo":
		return msg.Photo, nil
	case "video":
		return msg.Video, nil
	case "animation", "gif":
		return msg.Animation, nil
	default:
		return nil, fmt.Errorf("unsupported media type: %s", msg.Media().MediaType())
	}
}

// Updated processAlbum to return messages
func processAlbum(ctx m.Context, albumID string, dataChan chan *m.Message, keyboard *telebot.ReplyMarkup) ([]m.Message, error) {
	suggestChanId := ctx.Get("suggest_id").(int64)

	defer delete(pendingAlbums, albumID)
	defer close(dataChan)

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var messages []*m.Message
	for {
		select {
		case <-timeoutCtx.Done():
			if len(messages) == 0 {
				return nil, fmt.Errorf("no media collected from album: %s", albumID)
			}

			// Sort by message ID to maintain original order
			sort.Slice(messages, func(i, j int) bool {
				return messages[i].ID < messages[j].ID
			})

			// Create album
			album := make(telebot.Album, 0, len(messages))
			for _, msg := range messages {
				media, err := createMediaItem(msg)
				if err != nil {
					log.Printf("Skipping unsupported media: %v", err)
					continue
				}
				album = append(album, media)
			}

			if len(album) == 0 {
				return nil, fmt.Errorf("no valud media in album")
			}

			admins := ctx.Get("admins").([]string)

			// Set caption on the last item
			senderName := ctx.Sender().FirstName
			if slices.Contains(admins, ctx.Sender().Username) {
				album.SetCaption("[Caro est infirma](https://t.me/caroinfirma) ❤️ [Suggest a post](https://t.me/Caro_est_infirma_bot)")
			} else {
				caption := fmt.Sprintf("From %s \n\n [Caro est infirma](https://t.me/caroinfirma) ❤️ [Suggest a post](https://t.me/Caro_est_infirma_bot)", senderName)
				album.SetCaption(caption)
			}

			// Send album
			msgs, err := ctx.Bot().SendAlbum(&telebot.Chat{ID: suggestChanId}, album, &telebot.SendOptions{
				ParseMode: telebot.ModeMarkdown,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to send album: %w", err)
			}

			// Send keyboard as reply to first message
			if keyboard != nil {
				_, err = ctx.Bot().Send(
					&telebot.Chat{ID: suggestChanId},
					fmt.Sprintf("%d", len(msgs)),
					&telebot.SendOptions{
						ReplyMarkup: keyboard,
						ReplyTo:     &msgs[0],
					},
				)
				if err != nil {
					log.Printf("Failed to send keyboard: %v", err)
				}
			}
			return msgs, nil

		case msg := <-dataChan:
			messages = append(messages, msg)
		}
	}
}
