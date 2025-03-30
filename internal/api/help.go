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
/sq: Show queue table rows
forw message from our channel: Delete this post from channel and history
photo/video/gif (forward/direct): Save media in queue
`

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
photo/video/gif (forward/direct): Send media to admins
`

	if err := ctx.Send(ans); err != nil {
		return err
	}

	return nil
}
