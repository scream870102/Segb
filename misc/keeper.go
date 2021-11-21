package misc

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Keeper struct {
	Guild   string
	Channel string
	Offset  int
	Session *discordgo.Session
}

func (k *Keeper) AwakeBOT() {
	channels, _ := k.Session.GuildChannels(k.Guild)
	for _, channel := range channels {
		if channel.ID == k.Channel {
			k.Session.ChannelMessageSend(channel.ID, fmt.Sprintf("Awake bot at %s", time.Now().Local().String()))
			delay := time.Duration(k.Offset) * time.Minute
			time.AfterFunc(delay, k.AwakeBOT)
		}
	}
}
