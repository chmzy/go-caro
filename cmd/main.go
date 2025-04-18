package main

import (
	"context"
	a "go-caro/internal/api"
	"go-caro/internal/config"
	"go-caro/internal/events"
	hr "go-caro/internal/repository/history"
	qr "go-caro/internal/repository/queue"
	"go-caro/internal/service/flusher"
	hs "go-caro/internal/service/history"
	qs "go-caro/internal/service/queue"
	"go-caro/internal/service/sender"
	"go-caro/pkg/tg"
	"go-caro/pkg/tg/middleware"
	"go-caro/pkg/tg/model"
	"log"
	"time"

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
	mainChan.ID = tgCfg.ChannelID()

	flusherCfg, err := config.NewFlusherConfig()
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
	bot, err := tg.NewBot(tgCfg.Token())
	if err != nil {
		pgxpool.Close()
		log.Fatalf("failed to init bot: %s", err.Error())
	}
	queueSender := sender.NewSender(historyRepo, queueRepo, bot, tgCfg.ChannelID(), tgCfg.SuggestionsID(), tgCfg.PostPeriodSec(), tgCfg.RepostPeriodSec(), pgCfg.Timezone())
	queueSender.Start(context.Background(), 1*time.Minute)

	historyFlusher := flusher.NewFlusher(historyRepo)
	historyFlusher.Start(context.Background(), flusherCfg.FlushPeriodSec())

	bot.Handle(events.HelpCommand, middleware.WithValue("admins", tgCfg.Admins(), api.Help))
	bot.Handle(events.QueueCommand, middleware.FromAdmin(tgCfg.Admins(), api.Queue, model.NOOP))
	bot.Handle(events.HistoryCommand, middleware.FromAdmin(tgCfg.Admins(), api.History, model.NOOP))
	bot.Handle(&telebot.Btn{Unique: events.ApplyButton}, api.ApplyButtonEvent)
	bot.Handle(&telebot.Btn{Unique: events.RejectButton}, api.RejectButtonEvent)
	bot.Handle(&telebot.Btn{Unique: events.DeleteButton}, api.DeleteButtonEvent)
	bot.Handle(telebot.OnChannelPost, middleware.IsForwarded(middleware.DoBefore(queueSender.ShrinkSendPeriod, api.OnChannelPost), api.OnChannelPost))
	bot.Handle(telebot.OnMedia,
		middleware.WithValues(map[string]any{"chan_id": tgCfg.ChannelID(), "suggest_id": tgCfg.SuggestionsID(), "admins": tgCfg.Admins()},
			api.OnMedia))

	log.Println("Bot is running...")
	bot.Start()
}
