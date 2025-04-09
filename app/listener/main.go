/*
Interacting with Discord API
*/

package listener

import (
	"errors"
	"fmt"
	"strings"

	"github.com/darklab8/fl-darkbot/app/consoler"
	"github.com/darklab8/fl-darkbot/app/prometheuser"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/darklab8/go-utils/utils"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	dg, err := discordgo.New("Bot " + settings.Env.DiscorderBotToken)
	logus.Log.CheckFatal(err, "failed to init discord")

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(consolerHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	logus.Log.CheckFatal(err, "error opening connection,")
	defer dg.Close()

	logus.Log.Info("Bot is now running.  Press CTRL-C to exit.")
	utils.SleepAwaitCtrlC()
	logus.Log.Info("gracefully closed discord conn")
}

func allowedMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	botID := s.State.User.ID
	messageAuthorID := m.Author.ID
	botCreatorID := "370435997974134785"

	if !strings.HasPrefix(m.Content, fmt.Sprintf("%s ", settings.Env.ConsolerPrefix)) {
		return false
	}

	// TODO implement ability to allow info quering depending on some config option
	// if strings.HasPrefix(m.Content, fmt.Sprintf("%s info", settings.Env.ConsolerPrefix)) {
	// 	return true
	// }

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if messageAuthorID == botID {
		return false
	}

	// Bots should not command it
	if m.Author.Bot {
		return false
	}

	guild, guildErr := s.Guild(m.GuildID)

	// If not guild, then exit
	if guildErr != nil {
		return false
	}
	if m.Member == nil {
		return false
	}

	isBotController := false
	allowed_role := "bot_controller"
	gildMemberRoles, err2 := s.GuildRoles(m.GuildID)
	if err2 == nil {
		for _, PlayerRoleID := range m.Member.Roles {
			for _, GuildRole := range gildMemberRoles {
				if GuildRole.ID == PlayerRoleID {
					if GuildRole.Name == allowed_role {
						isBotController = true
					}
				}
			}
		}
	}

	var guild_name string = guild.Name

	// if message not from guild owner, bot creator or person with role bot_controller, then ignore
	if guild.OwnerID != messageAuthorID &&
		botCreatorID != messageAuthorID &&
		!isBotController {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("ERR access denied. You must be server owner or person with role named '%s' in order to command me!", allowed_role))
		prometheuser.ListenerIsAllowedOperations(guild_name, string(m.ChannelID), errors.New("not_allowed_by_listener")).Inc()
		return false
	}

	prometheuser.ListenerIsAllowedOperations(guild_name, string(m.ChannelID), nil).Inc()
	return true
}

var console *consoler.Consoler

func init() {
	console = consoler.NewConsoler(settings.Dbpath)
}

func Chunks(s string, chunkSize int) []string {
	if len(s) == 0 {
		return nil
	}
	if chunkSize >= len(s) {
		return []string{s}
	}
	var chunks []string = make([]string, 0, (len(s)-1)/chunkSize+1)
	currentLen := 0
	currentStart := 0
	for i := range s {
		if currentLen == chunkSize {
			chunks = append(chunks, s[currentStart:i])
			currentLen = 0
			currentStart = i
		}
		currentLen++
	}
	chunks = append(chunks, s[currentStart:])
	return chunks
}

func consolerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	is_allowed := allowedMessage(s, m)

	if !is_allowed {
		return
	}

	channelID := types.DiscordChannelID(m.ChannelID)
	rendered := console.Execute(m.Content, channelID)

	if rendered != "" {
		var err error
		chunked_strings := Chunks(rendered, 1950)
		for _, chunk_str := range chunked_strings {
			_, err := s.ChannelMessageSend(m.ChannelID, chunk_str)
			if err != nil {
				break
			}
		}
		if logus.Log.CheckWarn(err, "failed to send message of consoler") {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintln("failed to send message of consoler with reason: ", err.Error()))
		}
	}
	logus.Log.Debug("consolerHandler finished", logus.ChannelID(channelID))
}
