package main

import (
	"context"
	a "go-caro/internal/api"
	"go-caro/internal/config"
	hr "go-caro/internal/repository/history"
	qr "go-caro/internal/repository/queue"
	hs "go-caro/internal/service/history"
	qs "go-caro/internal/service/queue"
	"go-caro/pkg/tg"
	"log"

	"go-caro/pkg/tg/middleware"

	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/telebot.v4"
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

	api := a.NewAPI(historyService, queueService)
	b, err := tg.NewBot(tgCfg.Token())
	if err != nil {
		pgxpool.Close()
		log.Fatalf("failed to init bot: %s", err.Error())
	}
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

	// b.Bot.Handle("/test", func(ctx telebot.Context) error {
	// 	ctx.Send(ctx.Data())
	// 	for _, arg := range ctx.Args() {
	// 		ctx.Send(arg)
	// 	}
	// 	return nil
	// })

	log.Println("Bot is running...")
	b.Bot.Start()
}

// func main() {
// 	// Bot settings
// 	pref := telebot.Settings{
// 		Token:  os.Getenv("TOKEN"),
// 		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
// 	}

// 	// Create a new bot instance
// 	bot, err := telebot.NewBot(pref)
// 	if err != nil {
// 		log.Fatalf("Failed to create bot: %v", err)
// 	}

// 	// Handle the /help command
// 	bot.Handle("/help", func(c telebot.Context) error {
// 		// The message link
// 		messageLink := "https://t.me/c/265186136/1533"

// 		// Extract chat ID and message ID from the link
// 		// The link format is: https://t.me/c/<chat_id>/<message_id>
// 		parts := strings.Split(messageLink, "/")
// 		if len(parts) < 5 {
// 			return c.Send("Invalid message link format.")
// 		}
// 		fmt.Println(c.Message().ID)
// 		err := c.Send(fmt.Sprintf("%d", c.Message().ID))
// 		if err != nil {
// 			fmt.Println(err)
// 		}

// 		// chatIDPart := parts[4]                                   // Extract chat ID part (e.g., "265186136")
// 		// messageIDPart := parts[5]                                // Extract message ID part (e.g., "1533")
// 		// chatID, _ := strconv.ParseInt("-100"+chatIDPart, 10, 64) // Convert to full chat ID
// 		// messageID, _ := strconv.Atoi(messageIDPart)              // Convert to integer

// 		// // Forward the photo message to the user
// 		// bot.Forward(c.Sender(), &telebot.Message{
// 		// 	Chat: &telebot.Chat{ID: chatID}, ID: messageID,
// 		// })

// 		return nil
// 	})

// 	// Handle photo messages
// 	bot.Handle(telebot.OnPhoto, func(c telebot.Context) error {
// 		// Get the original photo and the message ID
// 		// photo := c.Message().Photo
// 		messageID := c.Message().ID

// 		// Generate the message link
// 		// userID := c.Sender().ID
// 		chatID := c.Chat().ID
// 		messageLink := fmt.Sprintf("https://t.me/c/%d/%d", chatID, messageID)

// 		// Forward the photo back to the user
// 		// c.Send(photo)

// 		// Forward the photo back to the user
// 		c.Send(fmt.Sprintf("Here's your photo! Original message: %s", messageLink))
// 		c.Send(telebot.Message{Caption: "Some capt"})

// 		return nil
// 	})

// 	// Start the bot
// 	log.Println("Bot is running...")
// 	bot.Start()
// }
