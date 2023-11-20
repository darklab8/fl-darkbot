/*
interface
- to capture message history from channel (internal)
- to create or replace message in channel (public?)
- delete message from channel (public?)
*/

package discorder

import (
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Discorder struct {
	dg *discordgo.Session
}

func NewClient() Discorder {
	d := Discorder{}
	dg, err := discordgo.New("Bot " + settings.Config.DiscorderBotToken)
	logus.CheckFatal(err, "failed to init discord")
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	d.dg = dg
	return d
}

func (d Discorder) SengMessage(channelID types.DiscordChannelID, content string) error {
	_, err := d.dg.ChannelMessageSend(string(channelID), content)
	logus.CheckWarn(err, "failed sending message in discorder", logus.ChannelID(channelID))
	return err
}

func (d Discorder) EditMessage(channelID types.DiscordChannelID, messageID types.DiscordMessageID, content string) error {
	_, err := d.dg.ChannelMessageEdit(string(channelID), string(messageID), content)
	logus.CheckWarn(err, "failed editing message in discorder", logus.ChannelID(channelID))
	return err
}

func (d Discorder) DeleteMessage(channelID types.DiscordChannelID, messageID types.DiscordMessageID) {
	err := d.dg.ChannelMessageDelete(string(channelID), string(messageID))
	logus.CheckWarn(err, "failed deleting message in discorder", logus.ChannelID(channelID))
}

type DiscordMessage struct {
	ID        types.DiscordMessageID
	Content   string
	Timestamp time.Time
}

func (d Discorder) GetLatestMessages(channelID types.DiscordChannelID) ([]DiscordMessage, error) {
	messagesLimitToGrab := 100 // max 100
	messages, err := d.dg.ChannelMessages(string(channelID), messagesLimitToGrab, "", "", "")
	if logus.CheckWarn(err, "Unable to get messages from channelId=", logus.ChannelID(channelID)) {
		return []DiscordMessage{}, err
	}

	result := []DiscordMessage{}

	for _, msg := range messages {
		result = append(result, DiscordMessage{
			ID:        types.DiscordMessageID(msg.ID),
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

	return result, nil
}

func (d Discorder) GetOwnerID(channelID types.DiscordChannelID) (types.DiscordOwnerID, error) {
	channel, err := d.dg.Channel(string(channelID))
	if logus.CheckError(err, "discord is not connected") {
		return types.DiscordOwnerID(""), err
	}
	channel_owner := types.DiscordOwnerID(channel.OwnerID)

	logus.Debug("channel.OwnerID=", logus.OwnerID(channel_owner))
	guildID := channel.GuildID

	guild, err := d.dg.Guild(guildID)
	if logus.CheckWarn(err, "unable to get Guild Owner", logus.ChannelID(channelID)) {
		return "", err
	}
	logus.CheckWarn(err, "Failed getting Guild ID")
	guild_owner_id := types.DiscordOwnerID(guild.OwnerID)
	logus.Debug("guild.OwnerID=", logus.OwnerID(guild_owner_id))

	return guild_owner_id, nil
}

type deduplicator struct {
	repeatCheckers []func(msgs []DiscordMessage) bool
}

func NewDeduplicator(checkers ...func(msgs []DiscordMessage) bool) *deduplicator {
	d := &deduplicator{
		repeatCheckers: checkers,
	}
	return d
}

func (d *deduplicator) isDuplicated(msgs []DiscordMessage) bool {
	for _, checker := range d.repeatCheckers {
		if checker(msgs) {
			return true
		}
	}
	return false
}

func (dg Discorder) SendDeduplicatedMsg(deduplicator *deduplicator, msg string, channel types.DiscordChannelID) {
	logus.Info("sent_message= " + msg)
	msgs, err := dg.GetLatestMessages(channel)

	if logus.CheckError(err, "failed to get discord latest msgs") {
		return
	}

	if deduplicator.isDuplicated(msgs) {
		logus.Debug("not sending duplicated", logus.ChannelID(channel), logus.DiscordMessage(msg))
		return
	}

	dg.SengMessage(channel, msg)
}
