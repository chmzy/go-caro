package api

import (
	"context"
	"fmt"
	m "go-caro/pkg/tg/model"
	"log"

	"gopkg.in/telebot.v4"
)

const (
	tgMsgLink = "https://t.me/c/%d/%s"
)

func (a *API) Queue(ctx m.Context) error {
	posts, err := a.queueService.Next(context.Background())
	if err != nil {
		return err
	}

	if len(posts) == 1 {
		msg := m.Post{
			ChatID:    posts[0].ChatID,
			MessageID: posts[0].MsgID,
		}

		opts := &telebot.SendOptions{
			ReplyTo: nil, // No reply needed
		}

		// _, err = ctx.Bot().Reply(msg, opts)
		err := ctx.Reply(msg, opts)
		if err != nil {
			return err
		}

		return nil
	}

	for _, post := range posts {
		log.Println(post)
		msg := fmt.Sprintf(tgMsgLink, post.ChatID, post.MsgID)

		opts := &telebot.SendOptions{
			ReplyTo: nil, // No reply needed
		}

		// _, err = ctx.Bot().Copy(&telebot.Chat{ID: tgCfg.ChannelID()}, msg, opts)
		err := ctx.Reply(msg, opts)
		if err != nil {
			return err
		}
	}

	return nil

}
