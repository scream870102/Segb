package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler struct{}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) Handle(s *discordgo.Session, e *discordgo.MessageCreate) {
	fmt.Printf("%s said : %s \n", e.Author.Username, e.Message.Content)
}
