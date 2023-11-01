package commands

import (
	"darkbot/configurator"
	"darkbot/configurator/models"
	"darkbot/consoler/commands/cmdgroup"
	"darkbot/consoler/helper"
	"darkbot/settings"
	"fmt"

	"github.com/spf13/cobra"
)

// Entrance into CLI
func createEntrance() *cobra.Command {
	command := &cobra.Command{
		Use:   "consoler",
		Short: "A brief description of your application",
		// Args:  cobra.MinimumNArgs(1),
		// // When commented out, HELP info is rendered
		// Run: func(cmd *cobra.Command, args []string) {
		// 	// Ignoring message and rendering nothing
		// },
		SilenceUsage:  true,
		SilenceErrors: true,
		Hidden:        true,
	}
	command.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	return command
}

func CreateConsoler(channelInfo helper.ChannelInfo) *cobra.Command {
	consolerCmd := createEntrance()

	rootGroup := cmdgroup.New(
		consolerCmd,
		channelInfo,
		cmdgroup.Command(settings.Config.ConsolerPrefix),
		cmdgroup.ShortDesc("Welcome to darkbot!"),
	)
	root := newRootCommands(&rootGroup)

	NewTagCommands(
		root.GetChild(
			root.CurrentCmd,
			cmdgroup.Command("base"),
			cmdgroup.ShortDesc("Base commands"),
		),
		configurator.NewConfiguratorBase(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	playerGroup := root.GetChild(
		root.CurrentCmd,
		cmdgroup.Command("player"),
		cmdgroup.ShortDesc("Player commands"),
	)
	NewTagCommands(
		playerGroup.GetChild(
			playerGroup.CurrentCmd,
			cmdgroup.Command("system"),
			cmdgroup.ShortDesc("System commands"),
		),
		configurator.NewConfiguratorSystem(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	NewTagCommands(
		playerGroup.GetChild(
			playerGroup.CurrentCmd,
			cmdgroup.Command("region"),
			cmdgroup.ShortDesc("Region commands"),
		),
		configurator.NewConfiguratorRegion(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	NewTagCommands(
		playerGroup.GetChild(
			playerGroup.CurrentCmd,
			cmdgroup.Command("friend"),
			cmdgroup.ShortDesc("Player friend commands"),
		),
		configurator.NewConfiguratorPlayerFriend(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	NewTagCommands(
		playerGroup.GetChild(
			playerGroup.CurrentCmd,
			cmdgroup.Command("enemy"),
			cmdgroup.ShortDesc("Player enemy commands"),
		),
		configurator.NewConfiguratorPlayerEnemy(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	alertGroup := root.GetChild(
		root.CurrentCmd,
		cmdgroup.Command("alert"),
		cmdgroup.ShortDesc("Alert commands"),
	)

	NewAlertBoolCommands[models.AlertBaseIfHealthDecreasing](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("base_health_is_decreasing"),
			cmdgroup.ShortDesc("Turn on to receive alert if base health is decreasing"),
		),
		configurator.NewCfgAlertBaseHealthIsDecreasing(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	NewAlertBoolCommands[models.AlertBaseIfUnderAttack](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("base_is_under_attack"),
			cmdgroup.ShortDesc("Turn on if base health is rapidly decreasing or attack declaration was declared"),
		),
		configurator.NewCfgAlertBaseIsUnderAttack(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	NewAlertThresholdCommands[models.AlertBaseHealthLowerThan](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("base_health_is_lower_than"),
			cmdgroup.ShortDesc("Set threshold of base health, below which you will receive alert"),
		),
		configurator.NewCfgAlertBaseHealthLowerThan(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	NewAlertThresholdCommands[models.AlertNeutralPlayersEqualOrGreater](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("player_neutral_count_above"),
			cmdgroup.ShortDesc("Set threshold, if above amount of neutral players will be preesent, you will receive alert"),
		),
		configurator.NewCfgAlertNeutralPlayersGreaterThan(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	NewAlertThresholdCommands[models.AlertEnemiesEqualOrGreater](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("player_enemy_count_above"),
			cmdgroup.ShortDesc("Set threshold, if above amount of enemy players will be preesent, you will receive alert"),
		),
		configurator.NewCfgAlertEnemyPlayersGreaterThan(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	NewAlertThresholdCommands[models.AlertFriendsEqualOrGreater](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("player_friend_count_above"),
			cmdgroup.ShortDesc("Set threshold, if above amount of friendly players will be preesent, you will receive alert"),
		),
		configurator.NewCfgAlertFriendPlayersGreaterThan(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	NewAlertSetStringCommand[models.AlertPingMessage](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("ping_message"),
			cmdgroup.ShortDesc("By default `<@DiscordServer.Owner.ID>`. You can change it to something else like `@here` or `@role`"),
		),
		configurator.NewCfgAlertPingMessage(configurator.NewConfigurator(channelInfo.Dbpath)),
	)

	return consolerCmd
}

type rootCommands struct {
	*cmdgroup.CmdGroup
	channels configurator.ConfiguratorChannel
}

func newRootCommands(
	cmdgroup *cmdgroup.CmdGroup,
) *rootCommands {
	r := &rootCommands{CmdGroup: cmdgroup}
	r.channels = configurator.NewConfiguratorChannel(r.Configurator)
	r.CreatePing()
	return r
}

func (r *rootCommands) CreatePing() {
	command := &cobra.Command{
		Use:   "ping",
		Short: "Check stuff is working",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ping called with args=", args)
			cmd.OutOrStdout().Write([]byte("Pong! from consoler"))
		},
	}
	r.CurrentCmd.AddCommand(command)
}

func (r *rootCommands) CreateConnect() {

	command := &cobra.Command{
		Use:   "connect",
		Short: "Connect bot to channel",
		Run: func(cmd *cobra.Command, args []string) {
			err := r.channels.Add(r.ChannelInfo.ChannelID).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR channel may be already connected, msg=" + err.Error()))
				return
			}
			cmd.OutOrStdout().Write([]byte("OK Channel is connected"))
		},
	}
	r.CurrentCmd.AddCommand(command)
}

func (r *rootCommands) CreateDisconnect() {
	command := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect bot from channel",
		Run: func(cmd *cobra.Command, args []string) {
			err := r.channels.Remove(r.ChannelInfo.ChannelID).GetError()
			if err != nil {
				cmd.OutOrStdout().Write([]byte("ERR channel may be already disconnected, msg=" + err.Error()))
				return
			}
			cmd.OutOrStdout().Write([]byte("OK Channel is disconnected"))
		},
	}
	r.CurrentCmd.AddCommand(command)
}
