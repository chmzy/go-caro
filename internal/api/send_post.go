package api

// import (
// 	"context"
// 	"fmt"
// 	"go-caro/internal/service/history/model"
// 	"log"
// 	"sort"
// 	"strconv"
// 	"time"

// 	"go-caro/pkg/tg"

// 	"gopkg.in/telebot.v4"
// )

// func (a *API) SendPost(bot *tg.TgBot) error {
// 	c := context.Background()
// 	// postingPeriod := ctx.Get("posting_period").(int64)
// 	// lastPost, err := a.historyService.GetLast(c)
// 	// if err != nil {
// 	// 	return fmt.Errorf("poller: historyRepo.getLast: %w\n", err)
// 	// }

// 	// if time.Now().UTC().Add(3*time.Hour).Sub(lastPost.PostedAt) < time.Duration(postingPeriod) {
// 	// 	return nil
// 	// }

// 	nextPosts, err := a.queueService.Next(c)
// 	if err != nil {
// 		return fmt.Errorf("poller: queueRepo.next: %w\n", err)
// 	}

// 	if len(nextPosts) == 0 {
// 		log.Println("No posts in queue")
// 		return fmt.Errorf("No posts in queue")
// 	}

// 	// Make newPost from queue
// 	var newPost []telebot.Editable
// 	for _, post := range nextPosts {
// 		newPost = append(newPost, telebot.StoredMessage{
// 			ChatID:    post.ChatID,
// 			MessageID: post.MsgID,
// 		})
// 	}
// 	// Sort messages bacause tg can return them unordered
// 	sort.Slice(newPost, func(i, j int) bool {
// 		msgId1, _ := newPost[i].MessageSig()
// 		msgId2, _ := newPost[j].MessageSig()
// 		return msgId1 <= msgId2
// 	})

// 	// Send media to main channel
// 	msgs, err := ctx.Bot().CopyMany(&telebot.Chat{ID: -1002040647793}, newPost)
// 	if err != nil {
// 		return fmt.Errorf("poller: copy many: %w\n", err)
// 	}

// 	// Add posts into history table after posting
// 	for _, m := range msgs {
// 		_, err := a.historyService.Create(context.Background(), &model.PostHistory{
// 			ID:       m.ID,
// 			AlbumID:  m.AlbumID,
// 			PostedAt: time.Now(),
// 		})
// 		if err != nil {
// 			log.Println("cant write to history: ", err.Error())
// 		}
// 	}

// 	// Delete posts from queue table after posting
// 	if len(nextPosts) > 1 {
// 		if err := a.queueService.DeleteByAlbumID(context.Background(), nextPosts[0].AlbumID); err != nil {
// 			log.Printf("poller: queueRepo: deleteBynewPostID: %s\n", err.Error())
// 		}
// 	} else {
// 		if err := a.queueService.DeleteByID(context.Background(), nextPosts[0].ID); err != nil {
// 			log.Printf("poller: queueRepo: deleteByID: %s\n", err.Error())
// 		}
// 	}

// 	if err := ctx.Bot().DeleteMany(newPost); err != nil {
// 		log.Println("Cant delete original posts from channel: ", err.Error())
// 	}

// 	// Delete keyboard if message from suggestion channel
// 	msID, chID := newPost[0].MessageSig()
// 	if chID == -1002504066662 {
// 		mid, _ := strconv.Atoi(msID)
// 		msg := telebot.StoredMessage{
// 			ChatID:    chID,
// 			MessageID: fmt.Sprintf("%d", mid+len(newPost)),
// 		}
// 		if err := ctx.Bot().Delete(msg); err != nil {
// 			log.Println("Cant delete msg with keyboard: ", err)
// 		}

// 	}

// 	log.Println("New post!")

// 	return nil

// }
