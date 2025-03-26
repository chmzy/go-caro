package poller

import (
	"context"
	"fmt"
	"go-caro/internal/repository"
	"go-caro/internal/service"
	"go-caro/internal/service/history/model"
	"go-caro/pkg/tg"
	"log"
	"sort"
	"strconv"
	"time"

	"gopkg.in/telebot.v4"
)

const (
	postingPeriod = 2 * time.Second
)

type poller struct {
	historyRepo service.HistoryService
	queueRepo   service.QueueService
	bot         *tg.TgBot
}

func NewPoller(historyRepo repository.HistoryRepository, queueRepo repository.QueueRepository, bot *tg.TgBot) service.PollerService {
	return &poller{
		historyRepo: historyRepo,
		queueRepo:   queueRepo,
		bot:         bot,
	}
}

func (p *poller) StartPolling(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// noop
			}
			lastPost, err := p.historyRepo.GetLast(ctx)
			if err != nil {
				log.Printf("poller: historyRepo.getLast: %s\n", err.Error())
				continue
			}

			if time.Now().UTC().Add(3*time.Hour).Sub(lastPost.PostedAt) < postingPeriod {
				continue
			}

			nextPosts, err := p.queueRepo.Next(ctx)
			if err != nil {
				log.Printf("poller: queueRepo.next: %s\n", err.Error())
			}

			if len(nextPosts) == 0 {
				log.Println("No posts in queue")
				continue
			}

			// Make newPost from queue
			var newPost []telebot.Editable
			for _, post := range nextPosts {
				newPost = append(newPost, telebot.StoredMessage{
					ChatID:    post.ChatID,
					MessageID: post.MsgID,
				})
			}
			// Sort messages bacause tg can return them unordered
			sort.Slice(newPost, func(i, j int) bool {
				msgId1, _ := newPost[i].MessageSig()
				msgId2, _ := newPost[j].MessageSig()
				return msgId1 <= msgId2
			})

			// Send media to main channel
			msgs, err := p.bot.CopyMany(&telebot.Chat{ID: -1002040647793}, newPost)
			if err != nil {
				log.Printf("poller: copy many: %s\n", err.Error())
				continue
			}

			// Add posts into history table after posting
			for _, m := range msgs {
				_, err := p.historyRepo.Create(context.Background(), &model.PostHistory{
					ID:       m.ID,
					AlbumID:  m.AlbumID,
					PostedAt: time.Now(),
				})
				if err != nil {
					log.Println("cant write to history: ", err.Error())
				}
			}

			// Delete posts from queue table after posting
			if len(nextPosts) > 1 {
				if err := p.queueRepo.DeleteByAlbumID(context.Background(), nextPosts[0].AlbumID); err != nil {
					log.Printf("poller: queueRepo: deleteBynewPostID: %s\n", err.Error())
				}
			} else {
				if err := p.queueRepo.DeleteByMsgID(context.Background(), nextPosts[0].MsgID); err != nil {
					log.Printf("poller: queueRepo: deleteByID: %s\n", err.Error())
				}
			}

			if err := p.bot.DeleteMany(newPost); err != nil {
				log.Println("Cant delete original posts from channel: ", err.Error())
			}

			// Delete keyboard if message from suggestion channel
			msID, chID := newPost[0].MessageSig()
			if chID == -1002504066662 {
				mid, _ := strconv.Atoi(msID)
				msg := telebot.StoredMessage{
					ChatID:    chID,
					MessageID: fmt.Sprintf("%d", mid+len(newPost)),
				}
				if err := p.bot.Delete(msg); err != nil {
					log.Println("Cant delete msg with keyboard: ", err)
				}

			}

			log.Println("New post!")

			continue
		}
	}()
}
