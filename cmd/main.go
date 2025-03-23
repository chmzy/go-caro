package main

import (
	"context"
	"fmt"
	a "go-caro/internal/api"
	"go-caro/internal/config"
	hr "go-caro/internal/repository/history"
	par "go-caro/internal/repository/pending_album"
	qr "go-caro/internal/repository/queue"
	hs "go-caro/internal/service/history"
	pas "go-caro/internal/service/pending_album"
	qs "go-caro/internal/service/queue"
	"go-caro/internal/service/queue/model"
	"go-caro/pkg/tg"
	"go-caro/pkg/tg/middleware"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/telebot.v4"
)

var (
	mainChan = &telebot.Chat{}
)

func main() {
	pgCfg, err := config.NewPGConfig()
	if err != nil {
		log.Fatal(err)
	}

	tgCfg, err := config.NewTGConfig()
	if err != nil {
		log.Fatal(err)
	}
	mainChan.ID = tgCfg.ID()

	pgxpool, err := pgxpool.New(context.Background(), pgCfg.DSN())
	if err != nil {
		log.Fatalf("failed to init pgxpool: %s", err.Error())
	}
	if err := pgxpool.Ping(context.Background()); err != nil {
		log.Fatalf("%s", err.Error())
	}

	historyRepo := hr.NewRepository(pgxpool)
	historyService := hs.NewService(historyRepo)

	queueRepo := qr.NewRepository(pgxpool)
	queueService := qs.NewService(queueRepo)

	pendigAlbumRepo := par.NewRepository(pgxpool)
	pendingAlbumService := pas.NewService(pendigAlbumRepo)

	api := a.NewAPI(historyService, queueService, pendingAlbumService)
	b, err := tg.NewBot(tgCfg.Token())
	if err != nil {
		pgxpool.Close()
		log.Fatalf("failed to init bot: %s", err.Error())
	}

	b.Bot.Handle(&telebot.Btn{Unique: "apply_action"}, func(ctx telebot.Context) error {
		// Extract user who clicked the button
		sender := ctx.Sender()

		// Check if message is media album
		if ctx.Message().ReplyTo != nil && ctx.Message().ReplyTo.AlbumID != "" {
			originalMessage := ctx.Message().ReplyTo
			albumLen, _ := strconv.Atoi(ctx.Message().Text)
			var albumMsgs []telebot.Editable
			for i := range albumLen {
				albumMsgs = append(albumMsgs, telebot.StoredMessage{ChatID: originalMessage.Chat.ID, MessageID: fmt.Sprintf("%d", originalMessage.ID+i)})

				_, err := queueRepo.Put(context.Background(), &model.PostQueue{
					Author:  sender.FirstName,
					AlbumID: originalMessage.AlbumID,
					MsgLink: model.Link{
						ChatID: originalMessage.Chat.ID,
						MsgID:  fmt.Sprintf("%d", originalMessage.ID+i),
					}})
				if err != nil {
					log.Println(err)
				}

			}
			// ctx.Bot().DeleteMany(albumMsgs)
			// ctx.Delete()
			inlineKeys := &telebot.ReplyMarkup{}
			btnApply := inlineKeys.Data("Will be posted soon...", "noop")
			btnReject := inlineKeys.Data("❌ Reject", "reject_action")
			inlineKeys.Inline(telebot.Row{btnApply, btnReject})

			ctx.Edit(ctx.Message(), &telebot.SendOptions{
				ReplyMarkup: inlineKeys,
			})

			return nil
		}

		_, err := queueRepo.Put(context.Background(), &model.PostQueue{
			Author:  sender.FirstName,
			AlbumID: "",
			MsgLink: model.Link{
				ChatID: ctx.Message().Chat.ID,
				MsgID:  fmt.Sprintf("%d", ctx.Message().ID),
			}})
		if err != nil {
			log.Println(err)
		}

		inlineKeys := &telebot.ReplyMarkup{}
		btnApproved := inlineKeys.Data("Will be posted soon...", "noop")
		btnDelete := inlineKeys.Data("❌ Delete from queue", "reject_action")
		inlineKeys.Inline(telebot.Row{btnApproved}, telebot.Row{btnDelete})

		err = ctx.Edit(inlineKeys)

		if err != nil {
			return err
		}

		// Edit the original message to reflect approval
		log.Printf("✅ Approved by %s", sender.FirstName)
		return nil
	})

	// b.Bot.Handle(&telebot.Btn{Unique: "suggest_reject"}, func(ctx telebot.Context) error {
	// 	// Extract user who clicked the button
	// 	user := ctx.Sender()
	// 	log.Printf("❌ Rejected by %s", user.FirstName)

	// 	// Edit the original message to reflect rejection
	// 	return nil
	// })

	b.Bot.Handle(telebot.OnChannelPost, api.OnChannelPost)

	b.Bot.Handle("/help", middleware.WithValue("admins", tgCfg.Admins(), api.Help))

	b.Bot.Handle(telebot.OnMedia,
		middleware.WithValue("chan_id", tgCfg.ID(),
			middleware.WithValue("admins", tgCfg.Admins(),
				api.OnMedia)))

	b.Bot.Handle("/queue", func(ctx telebot.Context) error {
		post, err := queueService.Next(context.Background())
		if err != nil {
			return err
		}

		msg := telebot.StoredMessage{
			ChatID:    post.MsgLink.ChatID,
			MessageID: post.MsgLink.MsgID,
		}

		opts := &telebot.SendOptions{
			ReplyTo: nil, // No reply needed
		}

		_, err = ctx.Bot().Copy(&telebot.Chat{ID: tgCfg.ID()}, msg, opts)
		if err != nil {
			return err
		}

		return nil
	})

	log.Println("Bot is running...")
	b.Bot.Start()
}

// func main() {
// 	tgCfg, err := config.NewTGConfig()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Initialize bot with long poller
// 	bot, err := telebot.NewBot(telebot.Settings{
// 		Token:  tgCfg.Token(),
// 		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
// 	})
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Handler for /copymany command
// 	bot.Handle("/copymany", func(c telebot.Context) error {
// 		// Get the replied message
// 		originalMsg := c.Message().ReplyTo
// 		if originalMsg == nil {
// 			return c.Send("❗ Please reply to a message to copy.")
// 		}

// 		// Copy the message to all recipients
// 		if _, err := bot.CopyMany(mainChan, []telebot.Editable{originalMsg}, &telebot.SendOptions{}); err != nil {
// 			log.Printf("Error copying message: %v", err)
// 			return c.Send("❌ Failed to copy message to some recipients.")
// 		}

// 		return c.Send("✅ Message copied successfully!")
// 	})

// 	log.Println("Bot started...")
// 	bot.Start()
// }
