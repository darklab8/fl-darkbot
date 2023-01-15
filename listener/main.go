/*
Interacting with Discord API
*/

package listener

import (
	"darkbot/consoler"
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

func consolerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	rendered := consoler.Consoler{}.New(m.Content).Execute(m.ChannelID).String()

	if rendered != "" {
		s.ChannelMessageSend(m.ChannelID, rendered)
	}
	fmt.Println("ChannelID=", m.ChannelID)
}
