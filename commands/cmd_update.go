package commands

import (
	"github.com/scream870102/segb/database"
)

type CmdUpdate struct{}

func (c *CmdUpdate) Invokes() []string {
	return []string{"update", "u"}
}

func (c *CmdUpdate) Description() string {
	return "update local cache from remote"
}

func (c *CmdUpdate) AdminRequired() bool {
	return true
}

func (c *CmdUpdate) Exec(ctx *Context) (err error) {
	ss := database.SpreadSheetInstance()
	ss.UpdateAllValueFromRemote()
	db := database.DatabaseInstance()
	db.Init()
	_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Update local cache from db successful")
	return err
}
