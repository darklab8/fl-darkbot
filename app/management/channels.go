/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package management

import (
	"encoding/json"
	"fmt"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/discorder"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"

	"github.com/spf13/cobra"
)

var channelDeleteCMD = &cobra.Command{
	Use:   "channel_delete",
	Short: "Delete channelID",
	Run: func(cmd *cobra.Command, args []string) {
		logus.Log.Info("Cmd is called with args=", logus.Args(args))
		channels := configurator.NewConfiguratorChannel(configurator.NewConfigurator(settings.Dbpath))

		fmt.Println("trying to delete ", len(args), args)
		var errors []error
		for _, channeID := range args {
			errors = append(errors, channels.Remove(types.DiscordChannelID(channeID)))
		}

		fmt.Println("error count", len(errors), " | ", errors)
	},
}
var channelListCMD = &cobra.Command{
	Use:   "channel_list",
	Short: "List channelIDs",
	Run: func(cmd *cobra.Command, args []string) {
		dg := discorder.NewClient()
		logus.Log.Info("Cmd is called with args=", logus.Args(args))
		channels := configurator.NewConfiguratorChannel(configurator.NewConfigurator(settings.Dbpath))

		channelIDs, _ := channels.List()
		fmt.Println("all channelIDs", len(channelIDs), " | ", channelIDs)

		var accessable_channels []types.DiscordChannelID
		var error_channels []types.DiscordChannelID
		var error_reason []error
		var channels_count_by_guild map[string]int = make(map[string]int)

		for _, channel := range channelIDs {
			_, err := dg.GetLatestMessages(channel)

			if ch, err := dg.GetDiscordSession().Channel(string(channel)); err == nil {
				if guild, err := dg.GetDiscordSession().Guild(ch.GuildID); err == nil {
					channels_count_by_guild[guild.Name] += 1
				}
			}

			if err == nil {
				accessable_channels = append(accessable_channels, channel)
			} else {
				error_channels = append(error_channels, channel)
				error_reason = append(error_reason, err)
			}
		}

		fmt.Println("accessable_channels", len(accessable_channels), " | ", accessable_channels)
		fmt.Println("error_channels", len(error_channels), " | ", error_channels)
		fmt.Println("error_reason", len(error_reason), " | ", error_reason)
		fmt.Println("guilds_count", len(channels_count_by_guild))

		data, _ := json.Marshal(channels_count_by_guild)
		fmt.Println("guilds_channel_counts", string(data))
	},
}

var channelInfoCMD = &cobra.Command{
	Use:   "channel_info",
	Short: "channel info",
	Run: func(cmd *cobra.Command, args []string) {
		logus.Log.Info("Cmd is called with args=", logus.Args(args))

		dis := discorder.NewClient().GetDiscordSession()

		channel_id := args[0]
		channel, err := dis.Channel(channel_id)

		if logus.Log.CheckError(err, "failed to get channel", logus.ChannelID(types.DiscordChannelID(channel_id))) {
			return
		}

		fmt.Println("channel=", channel)
		fmt.Println("channel.OwnerID=", channel.OwnerID)

		fmt.Println("channel.GuildID=", channel.GuildID)

		guild, err := dis.Guild(channel.GuildID)

		if logus.Log.CheckError(err, "failed to get guild="+channel.GuildID, logus.ChannelID(types.DiscordChannelID(channel_id))) {
			return
		}

		fmt.Println("guild=", guild)
		fmt.Println("guild.Name=", guild.Name)
		fmt.Println("guild.OwnerID=", guild.OwnerID)
	},
}

func init() {
	rootCmd.AddCommand(channelDeleteCMD)
	rootCmd.AddCommand(channelListCMD)
	rootCmd.AddCommand(channelInfoCMD)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// migrateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
