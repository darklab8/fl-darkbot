/*
TODO implement interface
- to capture message history from channel (internal)
- to create or replace message in channel (public?)
- delete message from channel (public?)
*/

package discorder

import (
	"darkbot/settings"
	"darkbot/utils"

	"github.com/bwmarrin/discordgo"
)

type Discorder struct {
	dg *discordgo.Session
}

func NewClient() Discorder {
	d := Discorder{}
	dg, err := discordgo.New("Bot " + settings.Config.DiscorderBotToken)
	utils.CheckPanic(err, "failed to init discord")
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	d.dg = dg
	return d
}

func (d Discorder) SengMessage(channelId string, message string) {
	d.dg.ChannelMessageSend(channelId, message)
}
