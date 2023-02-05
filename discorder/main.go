/*
interface
- to capture message history from channel (internal)
- to create or replace message in channel (public?)
- delete message from channel (public?)
*/

package discorder

import (
	"darkbot/settings"
	"darkbot/utils/logger"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Discorder struct {
	dg *discordgo.Session
}

func NewClient() Discorder {
	d := Discorder{}
	dg, err := discordgo.New("Bot " + settings.Config.DiscorderBotToken)
	logger.CheckPanic(err, "failed to init discord")
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	d.dg = dg
	return d
}

func (d Discorder) SengMessage(channelID string, content string) error {
	_, err := d.dg.ChannelMessageSend(channelID, content)
	logger.CheckWarn(err)
	return err
}

func (d Discorder) EditMessage(channelID string, messageID string, content string) error {
	_, err := d.dg.ChannelMessageEdit(channelID, messageID, content)
	logger.CheckWarn(err)
	return err
}

func (d Discorder) DeleteMessage(channelID string, messageID string) {
	err := d.dg.ChannelMessageDelete(channelID, messageID)
	logger.CheckWarn(err)
}

type DiscordMessage struct {
	ID        string
	Content   string
	Timestamp time.Time
}

func (d Discorder) GetLatestMessages(channelID string) []DiscordMessage {
	messagesLimitToGrab := 100 // max 100
	messages, err := d.dg.ChannelMessages(channelID, messagesLimitToGrab, "", "", "")
	if err != nil {
		logger.CheckWarn(err, "Unable to get messages from channelId=", channelID)
		return []DiscordMessage{}
	}

	result := []DiscordMessage{}

	for _, msg := range messages {
		result = append(result, DiscordMessage{
			ID:        msg.ID,
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
		})
	}

	// Just to be sure to have it deleted
	for index, _ := range messages {
		for index2, _ := range messages[index].Attachments {
			messages[index].Attachments[index2] = nil
		}
		for index2, _ := range messages[index].Embeds {
			messages[index].Embeds[index2] = nil
		}
		for index2, _ := range messages[index].MentionChannels {
			messages[index].MentionChannels[index2] = nil
		}
		for index2, _ := range messages[index].Mentions {
			messages[index].Mentions[index2] = nil
		}
		for index2, _ := range messages[index].Reactions {
			messages[index].Reactions[index2] = nil
		}
		for index2, _ := range messages[index].StickerItems {
			messages[index].StickerItems[index2] = nil
		}
		messages[index] = nil
	}
	messages = nil

	return result
}
