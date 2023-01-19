/*
Interacting with Discord API
*/

package listener

import (
	"darkbot/consoler"
	"darkbot/consoler/helper"
	"darkbot/settings"
	"darkbot/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	dg, err := discordgo.New("Bot " + settings.Config.DiscorderBotToken)
	utils.CheckPanic(err, "failed to init discord")

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(consolerHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	utils.CheckPanic(err, "error opening connection,")
	defer dg.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	utils.SleepAwaitCtrlC()
	fmt.Println("gracefully closed discord conn")
}

func allowedMessage(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	botID := s.State.User.ID
	messageAuthorID := m.Author.ID
	botCreatorID := "370435997974134785"

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

	var isBotController bool
	for _, roleID := range m.Member.Roles {
		role, err := s.State.Role(m.GuildID, roleID)
		if err != nil {
			continue
		}

		if role.Name == "bot_controller" {
			isBotController = true
		}
	}

	// if message not from guild owner, bot creator or person with role bot_controller, then ignore
	if guild.OwnerID != messageAuthorID &&
		botCreatorID != messageAuthorID &&
		!isBotController {
		return false
	}

	return true
}

func consolerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	if !allowedMessage(s, m) {
		return
	}

	rendered := consoler.Consoler{}.New(m.Content).Execute(helper.ChannelInfo{ChannelID: m.ChannelID}).String()

	if rendered != "" {
		s.ChannelMessageSend(m.ChannelID, rendered)
	}
	fmt.Println("ChannelID=", m.ChannelID)
}
