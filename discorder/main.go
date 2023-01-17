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

func (d Discorder) SengMessage(channelID string, content string) error {
	_, err := d.dg.ChannelMessageSend(channelID, content)
	utils.CheckWarn(err)
	return err
}

func (d Discorder) EditMessage(channelID string, messageID string, content string) error {
	_, err := d.dg.ChannelMessageEdit(channelID, messageID, content)
	utils.CheckWarn(err)
	return err
}

func (d Discorder) DeleteMessage(channelID string, messageID string) {
	d.dg.ChannelMessageDelete(channelID, messageID)
}

func (d Discorder) GetLatestMessages(channelID string) []*discordgo.Message {
	messagesLimitToGrab := 100 // max 100
	messages, err := d.dg.ChannelMessages(channelID, messagesLimitToGrab, "", "", "")
	if err != nil {
		utils.CheckWarn(err, "Unable to get messages from channelId=", channelID)
		return []*discordgo.Message{}
	}

	// Checking messages content
	// for _, msg := range messages {
	// 	fmt.Println(msg.Content)
	// }

	return messages
}
