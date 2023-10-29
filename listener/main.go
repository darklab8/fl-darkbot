/*
Interacting with Discord API
*/

package listener

import (
	"darkbot/consoler"
	"darkbot/consoler/helper"
	"darkbot/settings"
	"darkbot/settings/utils"
	"darkbot/settings/utils/logger"
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	dg, err := discordgo.New("Bot " + settings.Config.DiscorderBotToken)
	logger.CheckPanic(err, "failed to init discord")

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(consolerHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	logger.CheckPanic(err, "error opening connection,")
	defer dg.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	utils.SleepAwaitCtrlC()
	fmt.Println("gracefully closed discord conn")
}

func allowedMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	botID := s.State.User.ID
	messageAuthorID := m.Author.ID
	botCreatorID := "370435997974134785"

	if !strings.HasPrefix(m.Content, settings.Config.ConsolerPrefix) {
		return false
	}

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

	// if message not from guild owner, bot creator or person with role bot_controller, then ignore
	if guild.OwnerID != messageAuthorID &&
		botCreatorID != messageAuthorID &&
		!isBotController {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("ERR access denied. You must be server owner or person with role named '%s' in order to command me!", allowed_role))
		return false
	}

	return true
}

func consolerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if !allowedMessage(s, m) {
		return
	}

	rendered := consoler.Consoler{}.New(m.Content).Execute(helper.ChannelInfo{ChannelID: m.ChannelID, Dbpath: settings.Dbpath}).String()

	if rendered != "" {
		s.ChannelMessageSend(m.ChannelID, rendered)
	}
	fmt.Println("ChannelID=", m.ChannelID)
}
