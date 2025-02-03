package api

import (
	mw "go-caro/pkg/tg/middleware"
	m "go-caro/pkg/tg/model"
)

func (a *API) Help(ctx m.Context) error {
	adminUsers := ctx.Get("admins").([]string)

	return mw.FromAdmin(adminUsers, a.helpAdmin, a.helpUser)(ctx)
}

func (a *API) helpAdmin(ctx m.Context) error {
	ans := `
Commands:

/start: Start the bot
/help: Show command hints
/qdp: Delete single post from queue by queue id
/qdg: Delete group from queue by group id
/dp: Delete post from channel by msg_id
/sh: Show history table rows
/sq: Show queue table rows
/dh: Delete n rows from history table
forw message from our channel: Delete this post from channel and history
photo/video/gif (forward/direct): Save media in queue`

	if err := ctx.Send(ans); err != nil {
		return err
	}

	return nil
}

func (a *API) helpUser(ctx m.Context) error {
	ans := `
Commands:

/start: Start the bot
/help: Show command hints
photo/video/gif (forward/direct): Send media to admin`

	if err := ctx.Send(ans); err != nil {
		return err
	}

	return nil
}
