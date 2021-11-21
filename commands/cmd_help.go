package commands

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/scream870102/segb/misc"
)

type CmdHelp struct{}

func (c *CmdHelp) Invokes() []string {
	return []string{"help", "h"}
}

func (c *CmdHelp) Description() string {
	return "Show manual"
}

func (c *CmdHelp) AdminRequired() bool {
	return false
}

func (c *CmdHelp) Exec(ctx *Context) (err error) {
	file, err := os.Open("help.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	embed := misc.CreateBasicEmbed("Manual", true)
	embed.Description = fmt.Sprintf("Hey %s!!!\nHere is the manual for this BOT", ctx.Message.Author.Mention())
	img := discordgo.MessageEmbedImage{
		URL: "https://i.imgur.com/Mw51xtx.jpg",
	}
	embed.Image = &img

	title := ""
	content := ""
	isTitleFind := false
	for _, t := range text {
		if !strings.HasPrefix(t, "~") {
			if !isTitleFind {
				title = t
				isTitleFind = true
			} else {
				field := discordgo.MessageEmbedField{
					Name:  title,
					Value: content,
				}
				embed.Fields = append(embed.Fields, &field)
				title = t
				content = ""
			}
		} else {
			content += fmt.Sprintf("%s\n", t)
		}
	}
	_, err = ctx.Session.ChannelMessageSendEmbed(ctx.Message.ChannelID, embed)
	return err
}
