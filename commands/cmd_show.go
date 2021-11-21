package commands

import "github.com/scream870102/segb/database"

type CmdShow struct{}

func (c *CmdShow) Invokes() []string {
	return []string{"show", "s"}
}

func (c *CmdShow) Description() string {
	return "Show the content"
}

func (c *CmdShow) AdminRequired() bool {
	return false
}

func (c *CmdShow) Exec(ctx *Context) (err error) {
	db := database.DatabaseInstance()
	v := db.GetRawValue(ctx.Message.GuildID, ctx.Args)

	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, v.Content)
	return nil
}
