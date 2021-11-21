package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/scream870102/segb/database"
	"github.com/scream870102/segb/misc"
)

type CmdList struct{}
type listArg struct {
	Author string
	Tags   []string
}

func (c *CmdList) Invokes() []string {
	return []string{"list", "l"}
}

func (c *CmdList) Description() string {
	return "list the content"
}

func (c *CmdList) AdminRequired() bool {
	return false
}

func (c *CmdList) Exec(ctx *Context) (err error) {
	db := database.DatabaseInstance()
	listArg := listArg{}
	isAuthorFilter, isTagFilter := listArg.parse(ctx.Args)
	if isAuthorFilter || isTagFilter {
		if isAuthorFilter && isTagFilter {
			if id, exist := listArg.tryGetUserId(listArg.Author, ctx.Message.Mentions); exist {
				authorFilter := db.GetRawValuesByAuthor(ctx.Message.GuildID, id)
				tagFilter := db.GetRawValues(ctx.Message.GuildID, listArg.Tags)
				intersection := database.Intersection(authorFilter, tagFilter)
				if len(intersection) == 0 {
					_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "There is no any matched content.")
					return err
				} else {
					c.replyWithListEmbed(intersection, ctx)
					return err
				}
			}
		} else if isAuthorFilter || !isTagFilter {
			if id, exist := listArg.tryGetUserId(listArg.Author, ctx.Message.Mentions); exist {
				authorFilter := db.GetRawValuesByAuthor(ctx.Message.GuildID, id)
				if len(authorFilter) == 0 {
					_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "There is no any matched content.")
					return err
				} else {
					c.replyWithListEmbed(authorFilter, ctx)
					return err
				}
			}
		}
	}

	v := db.GetRawValues(ctx.Message.GuildID, ctx.Args)
	if len(v) == 0 {
		_, err = ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "There is no any matched content with tags.")
		return err
	}
	c.replyWithListEmbed(v, ctx)
	return err
}

func (c *CmdList) replyWithListEmbed(v []database.RawValue, ctx *Context) {
	totalPage := len(v)/misc.EmbedLimitField + 1
	for p := 0; p < totalPage; p++ {
		s := p * misc.EmbedLimitField
		e := s + misc.EmbedLimitField
		if e >= len(v) {
			e = len(v)
		}
		toShow := v[s:e]
		title := fmt.Sprintf("Page : %d/%d", p+1, totalPage)
		embed := c.createListContentEmbed(toShow, title, ctx)
		ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	}

}

func (c *CmdList) createListContentEmbed(v []database.RawValue, title string, ctx *Context) *discordgo.MessageEmbed {
	embed := misc.CreateBasicEmbed(title, true)
	image := discordgo.MessageEmbedImage{
		URL: "https://i.imgur.com/8LFgtgg.jpg",
	}
	embed.Image = &image
	embed.Fields = make([]*discordgo.MessageEmbedField, 0)
	for i, value := range v {
		tagValue := ""
		for _, tag := range value.Tags {
			tagValue += fmt.Sprintf("%s ,", tag)
		}
		tagField := discordgo.MessageEmbedField{
			Name: fmt.Sprintf("%d", i+1),
			Value: fmt.Sprintf("Tags : %s\n Author: %s \n Id : %s \n Content : %s \n",
				tagValue,
				ctx.Message.Author.Mention(),
				value.Id.String(),
				value.Content),
		}
		embed.Fields = append(embed.Fields, &tagField)
	}
	return embed
}

func (l *listArg) parse(args []string) (bool, bool) {
	isAuthorFind := false
	isTagFind := false
	for _, v := range args {
		if strings.Contains(v, "Author") || strings.Contains(v, "author") {
			reg, _ := regexp.Compile("(\\w+)=(.*)")
			res := reg.FindStringSubmatch(v)
			if len(res) > 0 {
				l.Author = res[2]
				isAuthorFind = true
			}
		} else {
			l.Tags = append(l.Tags, v)
			isTagFind = true
		}
	}
	return isAuthorFind, isTagFind
}

func (l *listArg) tryGetUserId(s string, users []*discordgo.User) (string, bool) {
	for _, user := range users {
		userM := user.Mention()
		userMExclamation := userM[:2] + "!" + userM[2:]
		if userM == s || userMExclamation == s {
			return user.ID, true
		}
	}
	return "", false
}
