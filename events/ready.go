package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type ReadyHandler struct{}

func NewReadyHandler() *ReadyHandler {
	return &ReadyHandler{}
}

func (h *ReadyHandler) Handler(dg *discordgo.Session, e *discordgo.Ready) {
	fmt.Println("BOT session is ready")
	fmt.Printf("Logged in as %s\n", e.User.String())
}
