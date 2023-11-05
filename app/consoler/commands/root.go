package commands

import (
	"darkbot/app/configurator"
	"darkbot/app/configurator/models"
	"darkbot/app/consoler/commands/cmdgroup"
	"darkbot/app/consoler/consoler_types"
	"darkbot/app/consoler/printer"
	"darkbot/app/settings"
	"darkbot/app/settings/logus"

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

func CreateConsoler(
	channelInfo *consoler_types.ChannelParams,
	configur *configurator.Configurator,
) *cobra.Command {
	consolerCmd := createEntrance()

	rootGroup := cmdgroup.NewCmdGroup(
		consolerCmd,
		channelInfo,
		configur,
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
		configurator.NewConfiguratorBase(configur),
		configurator.NewConfiguratorChannel(configur),
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
		configurator.NewConfiguratorSystem(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewTagCommands(
		playerGroup.GetChild(
			playerGroup.CurrentCmd,
			cmdgroup.Command("region"),
			cmdgroup.ShortDesc("Region commands"),
		),
		configurator.NewConfiguratorRegion(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewTagCommands(
		playerGroup.GetChild(
			playerGroup.CurrentCmd,
			cmdgroup.Command("friend"),
			cmdgroup.ShortDesc("Player friend commands"),
		),
		configurator.NewConfiguratorPlayerFriend(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewTagCommands(
		playerGroup.GetChild(
			playerGroup.CurrentCmd,
			cmdgroup.Command("enemy"),
			cmdgroup.ShortDesc("Player enemy commands"),
		),
		configurator.NewConfiguratorPlayerEnemy(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewTagCommands(
		playerGroup.GetChild(
			playerGroup.CurrentCmd,
			cmdgroup.Command("event"),
			cmdgroup.ShortDesc("Player event commands"),
		),
		configurator.NewConfiguratorPlayerEvent(configur),
		configurator.NewConfiguratorChannel(configur),
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
		configurator.NewCfgAlertBaseHealthIsDecreasing(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewAlertBoolCommands[models.AlertBaseIfUnderAttack](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("base_is_under_attack"),
			cmdgroup.ShortDesc("Turn on if base health is rapidly decreasing or attack declaration was declared"),
		),
		configurator.NewCfgAlertBaseIsUnderAttack(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewAlertThresholdCommands[models.AlertBaseHealthLowerThan](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("base_health_is_lower_than"),
			cmdgroup.ShortDesc("Set threshold of base health, below which you will receive alert"),
		),
		configurator.NewCfgAlertBaseHealthLowerThan(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewAlertThresholdCommands[models.AlertNeutralPlayersEqualOrGreater](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("player_neutral_count_above"),
			cmdgroup.ShortDesc("Set threshold, if above amount of neutral players will be preesent, you will receive alert"),
		),
		configurator.NewCfgAlertNeutralPlayersGreaterThan(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewAlertThresholdCommands[models.AlertEnemiesEqualOrGreater](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("player_enemy_count_above"),
			cmdgroup.ShortDesc("Set threshold, if above amount of enemy players will be preesent, you will receive alert"),
		),
		configurator.NewCfgAlertEnemyPlayersGreaterThan(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewAlertThresholdCommands[models.AlertFriendsEqualOrGreater](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("player_friend_count_above"),
			cmdgroup.ShortDesc("Set threshold, if above amount of friendly players will be preesent, you will receive alert"),
		),
		configurator.NewCfgAlertFriendPlayersGreaterThan(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewAlertSetStringCommand[models.AlertPingMessage](
		alertGroup.GetChild(
			alertGroup.CurrentCmd,
			cmdgroup.Command("ping_message"),
			cmdgroup.ShortDesc("By default `<@DiscordServer.Owner.ID>`. You can change it to something else like `@here` or `@role`"),
		),
		configurator.NewCfgAlertPingMessage(configur),
		configurator.NewConfiguratorChannel(configur),
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
	r.CreateConnect()
	r.CreateDisconnect()
	return r
}

func (r *rootCommands) CreatePing() {
	command := &cobra.Command{
		Use:   "ping",
		Short: "Check stuff is working",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Debug("ping called with args=", logus.Args(args))
			printer.Println(cmd, "Pong! from consoler")
		},
	}
	r.CurrentCmd.AddCommand(command)
}

func (r *rootCommands) CreateConnect() {

	command := &cobra.Command{
		Use:   "connect",
		Short: "Connect bot to channel",
		Run: func(cmd *cobra.Command, args []string) {
			err := r.channels.Add(r.GetChannelID())
			if err != nil {
				printer.Println(cmd, "ERR channel may be already connected, msg="+err.Error())
				return
			}
			printer.Println(cmd, "OK Channel is connected")
		},
	}
	r.CurrentCmd.AddCommand(command)
}

func (r *rootCommands) CreateDisconnect() {
	command := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect bot from channel",
		Run: func(cmd *cobra.Command, args []string) {
			err := r.channels.Remove(r.GetChannelID())
			if err != nil {
				printer.Println(cmd, "ERR channel may be already disconnected, msg="+err.Error())
				return
			}
			printer.Println(cmd, "OK Channel is disconnected")
		},
	}
	r.CurrentCmd.AddCommand(command)
}
