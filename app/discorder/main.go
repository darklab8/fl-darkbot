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
	"darkbot/app/settings/utils"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Discorder struct {
	dg *discordgo.Session
}

func (d *Discorder) GetDiscordSession() *discordgo.Session {
	return d.dg
}

func NewClient() *Discorder {
	d := &Discorder{}
	dg, err := discordgo.New("Bot " + settings.Config.DiscorderBotToken)
	logus.CheckFatal(err, "failed to init discord")
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	d.dg = dg
	return d
}

func (d *Discorder) SengMessage(channelID types.DiscordChannelID, content string) error {
	_, err := d.dg.ChannelMessageSend(string(channelID), content)
	logus.CheckWarn(err, "failed sending message in discorder", logus.ChannelID(channelID))
	return err
}

func (d *Discorder) EditMessage(channelID types.DiscordChannelID, messageID types.DiscordMessageID, content string) error {
	var err error
	utils.TimeMeasure(func() {
		msg, err := d.dg.ChannelMessageEdit(string(channelID), string(messageID), content)
		logus.CheckWarn(err, "failed editing message in discorder", logus.ChannelID(channelID))
		logus.Debug(fmt.Sprintf("Discorder.EditMessage.msg=%v", msg))
	}, fmt.Sprintf("Discorder.EditMessage content=%s", content), logus.ChannelID(channelID), logus.MessageID(messageID))
	return err
}

func (d *Discorder) DeleteMessage(channelID types.DiscordChannelID, messageID types.DiscordMessageID) error {
	err := d.dg.ChannelMessageDelete(string(channelID), string(messageID))
	logus.CheckWarn(err, "failed deleting message in discorder", logus.ChannelID(channelID))
	return err
}

type DiscordMessage struct {
	ID        types.DiscordMessageID
	Content   string
	Timestamp time.Time
	Embeds    []*discordgo.MessageEmbed
}

func (d *Discorder) GetLatestMessages(channelID types.DiscordChannelID) ([]*DiscordMessage, error) {
	messagesLimitToGrab := 100 // max 100
	messages, err := d.dg.ChannelMessages(string(channelID), messagesLimitToGrab, "", "", "")
	if logus.CheckWarn(err, "Unable to get messages from channelId=", logus.ChannelID(channelID)) {
		return []*DiscordMessage{}, err
	}

	result := []*DiscordMessage{}

	for _, msg := range messages {
		result = append(result, &DiscordMessage{
			ID:        types.DiscordMessageID(msg.ID),
			Content:   msg.Content,
			Timestamp: msg.Timestamp,
			Embeds:    msg.Embeds,
		})
	}

	return result, nil
}

func (d *Discorder) GetOwnerID(channelID types.DiscordChannelID) (types.DiscordOwnerID, error) {
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

type Deduplicator struct {
	dupCheckers []func() bool
}

type DuplicatedError struct {
}

func (d DuplicatedError) Error() string { return "This msg is duplicated" }

func NewDeduplicator(isDuplicaters ...func() bool) *Deduplicator {
	d := &Deduplicator{
		dupCheckers: isDuplicaters,
	}
	return d
}

func (d *Deduplicator) isDuplicated() bool {
	for _, isDup := range d.dupCheckers {
		if isDup() {
			return true
		}
	}
	return false
}

func (d *Discorder) SendDeduplicatedMsg(
	deduplicator *Deduplicator,
	channel types.DiscordChannelID,
	send_callback func(channel types.DiscordChannelID, dg *discordgo.Session) error,
) error {
	if deduplicator.isDuplicated() {
		logus.Debug("not sending duplicated", logus.ChannelID(channel))
		return DuplicatedError{}
	}

	err := send_callback(channel, d.dg)
	logus.CheckWarn(err, "failed sending message in discorder", logus.ChannelID(channel))
	return err
}
