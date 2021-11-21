package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"github.com/scream870102/segb/database"
	"github.com/scream870102/segb/misc"
)

type CmdAdd struct{}

func (c *CmdAdd) Invokes() []string {
	return []string{"add", "a"}
}

func (c *CmdAdd) Description() string {
	return "add new content"
}

func (c *CmdAdd) AdminRequired() bool {
	return false
}

func (c *CmdAdd) Exec(ctx *Context) (err error) {
	content := ctx.Args[len(ctx.Args)-1]
	tags := ctx.Args[0 : len(ctx.Args)-1]
	ss := database.SpreadSheetInstance()
	id, _ := uuid.NewUUID()
	rawValue := database.RawValue{
		Id:      id,
		Content: content,
		Author:  ctx.Message.Author.ID,
		Tags:    tags,
	}
	if ss.TryAdd(rawValue, ctx.Message.GuildID) {
		embed := misc.CreateBasicEmbed("Add new content successful", true)
		thumbnail := discordgo.MessageEmbedThumbnail{
			URL:    ctx.Message.Author.AvatarURL("128"),
			Width:  128,
			Height: 128,
		}
		embed.Thumbnail = &thumbnail

		embed.Description = fmt.Sprintf("Congratulation @!%s!!!\nYou are the best momoko oshi", ctx.Message.Author.Mention())

		idField := discordgo.MessageEmbedField{
			Name: "Id", Value: id.String(),
		}

		tagValue := ""
		for _, tag := range tags {
			tagValue += fmt.Sprintf("%s ,", tag)
		}
		tagField := discordgo.MessageEmbedField{
			Name: "Tag", Value: tagValue,
		}

		contentField := discordgo.MessageEmbedField{
			Name: "Content", Value: content,
		}
		embed.Fields = []*discordgo.MessageEmbedField{&idField, &tagField, &contentField}

		_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)

	} else {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, fmt.Sprintf("Add new content failed %s", err.Error()))
	}
	return err
}


