package misc

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

const (
	EmbedLimitTitle       = 256
	EmbedLimitDescription = 2048
	EmbedLimitFieldValue  = 1024
	EmbedLimitFieldName   = 256
	EmbedLimitField       = 25
	EmbedLimitFooter      = 2048
	EmbedLimit            = 4000
)

func CreateBasicEmbed(title string, isWithFooter bool) *discordgo.MessageEmbed {
	result := discordgo.MessageEmbed{}
	result.Title = title
	result.Color = 0xf05d9a
	result.Timestamp = time.Now().Format(time.RFC3339)
	if isWithFooter {
		footer := discordgo.MessageEmbedFooter{}
		footer.IconURL = "https://i.imgur.com/oC0Yo5s.png"
		footer.Text = "Created by scream870102"
		result.Footer = &footer
	}
	return &result
}
