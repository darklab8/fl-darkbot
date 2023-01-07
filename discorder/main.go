/*
Interacting with Discord API
*/

package discorder

import (
	"darkbot/consoler"
	"darkbot/settings"
	"darkbot/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Run() {
	dg, err := discordgo.New("Bot " + settings.Config.Discorder.Discord.Bot.Token)
	utils.CheckPanic(err, "failed to init discord")

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(pingHandler)
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

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func pingHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

func consolerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	rendered := consoler.Execute(m.Content)

	if rendered != "" {
		s.ChannelMessageSend(m.ChannelID, rendered)
	}
}
