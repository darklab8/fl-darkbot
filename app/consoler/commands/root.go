package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/configurator/models"
	"github.com/darklab8/fl-darkbot/app/consoler/commands/cmdgroup"
	"github.com/darklab8/fl-darkbot/app/consoler/consoler_types"
	"github.com/darklab8/fl-darkbot/app/consoler/printer"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"

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
		cmdgroup.Command(settings.Env.ConsolerPrefix),
		cmdgroup.ShortDesc("Welcome to darkbot!"),
	)
	root := newRootCommands(&rootGroup, configur)

	baseGroup := root.GetChild(
		root.CurrentCmd,
		cmdgroup.Command("base"),
		cmdgroup.ShortDesc("Base commands"),
	)

	NewTagCommands(
		baseGroup.GetChild(
			baseGroup.CurrentCmd,
			cmdgroup.Command("tags"),
			cmdgroup.ShortDesc("base tags u add for tracking"),
		),
		configurator.NewConfiguratorBase(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	NewAlertSetStringCommand[models.ConfigBaseOrderingKey](
		baseGroup.GetChild(
			baseGroup.CurrentCmd,
			cmdgroup.Command("order_by"),
			cmdgroup.ShortDesc(fmt.Sprintf("changing ordering to one of allowed keys: %v", models.ConfigBaseOrderingKeyAllowedTags)),
		),
		configurator.NewCfgBaseOrderingKey(configur),
		configurator.NewConfiguratorChannel(configur),
		models.ConfigBaseOrderingKeyAllowedTags,
	)

	forumGroup := root.GetChild(
		root.CurrentCmd,
		cmdgroup.Command("forum"),
		cmdgroup.ShortDesc("forum commands"),
	)
	forumThreadGroup := forumGroup.GetChild(
		forumGroup.CurrentCmd,
		cmdgroup.Command("thread"),
		cmdgroup.ShortDesc("track by thread name"),
	)
	NewTagCommands(
		forumThreadGroup.GetChild(
			forumThreadGroup.CurrentCmd,
			cmdgroup.Command("watch"),
			cmdgroup.ShortDesc("Watch commands"),
		),
		configurator.NewConfiguratorForumWatch(configur),
		configurator.NewConfiguratorChannel(configur),
	)
	NewTagCommands(
		forumThreadGroup.GetChild(
			forumThreadGroup.CurrentCmd,
			cmdgroup.Command("ignore"),
			cmdgroup.ShortDesc("Ignore commands"),
		),
		configurator.NewConfiguratorForumIgnore(configur),
		configurator.NewConfiguratorChannel(configur),
	)

	forumSubforumGroup := forumGroup.GetChild(
		forumGroup.CurrentCmd,
		cmdgroup.Command("subforum"),
		cmdgroup.ShortDesc("track by subforum name"),
	)
	NewTagCommands(
		forumSubforumGroup.GetChild(
			forumSubforumGroup.CurrentCmd,
			cmdgroup.Command("watch"),
			cmdgroup.ShortDesc("Watch commands"),
		),
		configurator.NewConfiguratorSubForumWatch(configur),
		configurator.NewConfiguratorChannel(configur),
	)
	NewTagCommands(
		forumSubforumGroup.GetChild(
			forumSubforumGroup.CurrentCmd,
			cmdgroup.Command("ignore"),
			cmdgroup.ShortDesc("Ignore commands"),
		),
		configurator.NewConfiguratorSubForumIgnore(configur),
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
		[]types.OrderKey{},
	)

	return consolerCmd
}

type rootCommands struct {
	*cmdgroup.CmdGroup
	channels configurator.ConfiguratorChannel
	*configurator.Configurators
}

func newRootCommands(
	cmdgroup *cmdgroup.CmdGroup,
	configur *configurator.Configurator,
) *rootCommands {
	r := &rootCommands{
		CmdGroup:      cmdgroup,
		Configurators: configurator.NewConfiguratorsFromConfigur(configur),
	}
	r.channels = configurator.NewConfiguratorChannel(r.Configurator)

	r.CreatePing()
	r.CreateConnect()
	r.CreateDisconnect()
	r.CreateConfig()

	NewInfoCommands(
		r.CmdGroup,
		r.channels,
	)
	return r
}

func (r *rootCommands) CreatePing() {
	command := &cobra.Command{
		Use:   "ping",
		Short: "Check stuff is working",
		Run: func(cmd *cobra.Command, args []string) {
			logus.Log.Debug("ping called with args=", logus.Args(args))
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

func (r *rootCommands) CreateConfig() {
	command := &cobra.Command{
		Use:   "conf",
		Short: "See all your configs",
		Run: func(cmd *cobra.Command, args []string) {
			var sb strings.Builder

			channel_id := r.GetChannelID()

			is_enabled_channel, _ := r.channels.IsEnabled(channel_id)

			sb.WriteString(fmt.Sprintln("darkbot version: ", os.Getenv("BUILD_VERSION")))
			sb.WriteString(fmt.Sprintln("is channel ", channel_id, "connected = ", strconv.FormatBool(is_enabled_channel)))

			if !is_enabled_channel {
				printer.Println(cmd, sb.String())
				return
			}

			// bases
			sb.WriteString("Bases:\n```\n")
			sb.WriteString(fmt.Sprintf("base tags = %#v\n", PrintList(r.Bases.Tags.TagsList2(channel_id))))
			sb.WriteString(fmt.Sprintf("base order_by = %#v\n", GetStatus(r.Configurators.Bases.OrderBy, channel_id)))
			sb.WriteString("\n```\n")

			// players
			sb.WriteString("Players:\n```\n")
			sb.WriteString(fmt.Sprintf("player region = %#v\n", PrintList(r.Players.Regions.TagsList2(channel_id))))
			sb.WriteString(fmt.Sprintf("player system = %#v\n", PrintList(r.Players.Systems.TagsList2(channel_id))))
			sb.WriteString(fmt.Sprintf("player friend = %#v\n", PrintList(r.Players.Friends.TagsList2(channel_id))))
			sb.WriteString(fmt.Sprintf("player enemy = %#v\n", PrintList(r.Players.Enemies.TagsList2(channel_id))))
			sb.WriteString(fmt.Sprintf("player event = %#v\n", PrintList(r.Players.Events.TagsList2(channel_id))))
			sb.WriteString("\n```\n")

			// forum
			sb.WriteString("Forum:\n```\n")
			sb.WriteString(fmt.Sprintf("forum subforum watch = %#v\n", PrintList(r.Forum.Subforum.Watch.TagsList2(channel_id))))
			sb.WriteString(fmt.Sprintf("forum subforum ignore = %#v\n", PrintList(r.Forum.Subforum.Ignore.TagsList2(channel_id))))
			sb.WriteString(fmt.Sprintf("forum thread watch = %#v\n", PrintList(r.Forum.Thread.Watch.TagsList2(channel_id))))
			sb.WriteString(fmt.Sprintf("forum thread ignore = %#v\n", PrintList(r.Forum.Thread.Ignore.TagsList2(channel_id))))
			sb.WriteString("\n```\n")

			// alerts
			sb.WriteString("Alerts:\n```\n")

			sb.WriteString(fmt.Sprintf("alert base_health_is_decreasing = %#v\n", GetStatus(r.Alerts.BaseHealthIsDecreasing, channel_id)))
			sb.WriteString(fmt.Sprintf("alert base_health_is_lower_than = %#v\n", GetStatus(r.Alerts.BaseHealthLowerThan, channel_id)))
			sb.WriteString(fmt.Sprintf("alert base_is_under_attack = %#v\n", GetStatus(r.Alerts.BaseIsUnderAttack, channel_id)))
			sb.WriteString(fmt.Sprintf("alert player_enemy_count_above = %#v\n", GetStatus(r.Alerts.EnemiesGreaterThan, channel_id)))
			sb.WriteString(fmt.Sprintf("alert player_friend_count_above = %#v\n", GetStatus(r.Alerts.FriendsGreaterThan, channel_id)))
			sb.WriteString(fmt.Sprintf("alert player_neutral_count_above = %#v\n", GetStatus(r.Alerts.NeutralsGreaterThan, channel_id)))

			value, err := r.Alerts.PingMessage.Status(channel_id)
			if err != nil {
				switch err.(type) {
				case configurator.ErrorZeroAffectedRows:
					sb.WriteString(fmt.Sprintln("alert ping_message = Server Owner"))
				default:
					sb.WriteString(fmt.Sprintln("alert ping_message = ", err.Error()))
				}
			} else {
				sb.WriteString(fmt.Sprintf("alert ping_message = %#v\n", fmt.Sprintf("%s", value)))
			}
			sb.WriteString("\n```\n")

			// return all
			printer.Println(cmd, sb.String())
		},
	}
	r.CurrentCmd.AddCommand(command)
}

type ConfStatus[T any] interface {
	Status(channelID types.DiscordChannelID) (T, error)
}

func PrintList[T any](smth []T) string {
	var sb strings.Builder

	sb.WriteString("[ ")
	for index, obj := range smth {
		sb.WriteString(fmt.Sprintf("%v", obj))

		if index != len(smth)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(" ]")
	return sb.String()
}

func GetStatus[T any](r ConfStatus[T], channelID types.DiscordChannelID) string {
	value, err := r.Status(channelID)

	if err != nil {
		switch err.(type) {
		case configurator.ErrorZeroAffectedRows:
			return "not set"
		default:
			return err.Error()
		}
	}

	return fmt.Sprint(value)
}
