package sender

import (
	"context"
	"fmt"
	"go-caro/internal/repository"
	"go-caro/pkg/tg"
	"log"
	"sort"
	"strconv"
	"time"

	modelserv "go-caro/internal/service/history/model"

	"gopkg.in/telebot.v4"
)

var (
	defaultSendPeriod            time.Duration
	defaultSendPeriodAfterRepost time.Duration
)

type sender struct {
	historyRepo   repository.HistoryRepository
	queueRepo     repository.QueueRepository
	bot           *tg.TgBot
	sendPeriod    time.Duration
	timezone      time.Duration
	mainChanID    int64
	suggestChanID int64
}

func NewSender(historyRepo repository.HistoryRepository, queueRepo repository.QueueRepository, bot *tg.TgBot, mainChanID, suggestChanID int64, sendPeriod, afterRepostPeriod, timezone int64) *sender {
	defaultSendPeriod = time.Duration(sendPeriod) * time.Second
	defaultSendPeriodAfterRepost = time.Duration(afterRepostPeriod) * time.Second

	return &sender{
		historyRepo:   historyRepo,
		queueRepo:     queueRepo,
		bot:           bot,
		sendPeriod:    defaultSendPeriod,
		mainChanID:    mainChanID,
		suggestChanID: suggestChanID,
		timezone:      time.Duration(timezone) * time.Hour,
	}
}

func (s *sender) Start(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// noop
			}
			lastPost, err := s.historyRepo.GetLast(ctx)
			if err != nil {
				log.Printf("sender: historyRepo.getLast: %s\n", err.Error())
				continue
			}

			if time.Now().UTC().Add(s.timezone).Sub(lastPost.PostedAt) < s.sendPeriod {
				log.Println("sender: next post in: ", s.sendPeriod-time.Now().UTC().Add(s.timezone).Sub(lastPost.PostedAt))
				continue
			}

			nextPosts, err := s.queueRepo.Next(ctx)
			if err != nil {
				log.Printf("sender:  %s\n", err.Error())
			}

			if len(nextPosts) == 0 {
				log.Println("No posts in queue")
				continue
			}

			// Make newPost from queue
			var newPosts []telebot.Editable
			for _, post := range nextPosts {
				newPosts = append(newPosts, telebot.StoredMessage{
					ChatID:    post.ChatID,
					MessageID: post.MsgID,
				})
			}
			// Sort messages bacause tg can return them unordered
			sort.Slice(newPosts, func(i, j int) bool {
				msgId1, _ := newPosts[i].MessageSig()
				msgId2, _ := newPosts[j].MessageSig()
				return msgId1 < msgId2
			})

			// Send media to main channel
			msgs, err := s.bot.CopyMany(&telebot.Chat{ID: s.mainChanID}, newPosts)
			if err != nil {
				log.Printf("sender: copy many: %s\n", err.Error())
				continue
			}

			// Add posts into history table after posting
			for _, m := range msgs {
				// TODO: CopyMany don't return *Chat object, so we cant extract ChatID from *telebot.Message struct
				// Need to do it manually
				_, err := s.historyRepo.Create(context.Background(), &modelserv.PostHistory{
					ID:       0,
					AlbumID:  nextPosts[0].AlbumID,
					ChatID:   s.mainChanID,
					MsgID:    fmt.Sprintf("%d", m.ID),
					PostedAt: time.Now().UTC().Add(s.timezone),
				})
				if err != nil {
					// Delete new posts if historyRepo failure occured
					log.Println("sender: ", err.Error())
					for _, msg := range msgs {
						if err := s.bot.Delete(telebot.StoredMessage{MessageID: fmt.Sprintf("%d", msg.ID), ChatID: msg.Chat.ID}); err != nil {
							log.Println("poller: delete post after history failure: ", err.Error())
						}
					}
					continue
				}
			}

			// Delete posts from queue table after posting
			if len(nextPosts) > 1 {
				if err := s.queueRepo.DeleteByAlbumID(context.Background(), nextPosts[0].AlbumID); err != nil {
					log.Printf("sender:  %s\n", err.Error())
				}
			} else {
				if err := s.queueRepo.DeleteByMsgID(context.Background(), nextPosts[0].MsgID); err != nil {
					log.Printf("sender: %s\n", err.Error())
				}
			}

			if err := s.bot.DeleteMany(newPosts); err != nil {
				log.Println("sender: cant delete original posts from channel: ", err.Error())
			}

			// Delete keyboard if message from suggestion channel
			if nextPosts[0].AlbumID != "" {
				msID, chID := newPosts[0].MessageSig()
				if chID == s.suggestChanID {
					msgIdInt, _ := strconv.Atoi(msID)
					msg := telebot.StoredMessage{
						ChatID:    chID,
						MessageID: fmt.Sprintf("%d", msgIdInt+len(nextPosts)),
					}
					if err := s.bot.Delete(msg); err != nil {
						log.Println("sender: cant delete keyboard: ", err)
					}

				}
			}
			if s.sendPeriod == defaultSendPeriodAfterRepost {
				s.RestoreSendPeriod()
			}

			log.Println("New post!")

			continue
		}
	}()
}

func (p *sender) ShrinkSendPeriod() {
	p.sendPeriod = defaultSendPeriodAfterRepost
}

func (p *sender) RestoreSendPeriod() {
	p.sendPeriod = defaultSendPeriod
}
