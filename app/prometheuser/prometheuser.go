package prometheuser

import (
	"net/http"
	"time"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/discorder"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// to add later
	// requestsPerChannel = promauto.NewCounterVec(prometheus.CounterOpts{
	// 	Name: "darkbot_channel_commands_called_count",
	// 	Help: "The total requests per channel",
	// }, []string{"guild_name", "channel_id"})
	// healthyChannels = promauto.NewCounterVec(prometheus.CounterOpts{
	// 	Name: "darkbot_health_count",
	// 	Help: "Healthy channel ping",
	// }, []string{"guild_name", "channel_id"})
	// erroredChannels = promauto.NewCounterVec(prometheus.CounterOpts{
	// 	Name: "darkbot_errors_count",
	// 	Help: "Channel ping resulted in error",
	// }, []string{"guild_name", "channel_id"})
	channelsPerGuilds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkbot_guilds_channels_gauge",
		Help: "The total number of channels used per guild",
	}, []string{"guild_name", "channel_id"})
)

//	func CommandsCalled(GuildID string, ChannelID string) prometheus.Counter {
//		return requestsPerChannel.WithLabelValues(GuildID, ChannelID)
//	}
//
//	func HealthyChannelCalled(GuildID string, ChannelID string) prometheus.Counter {
//		return healthyChannels.WithLabelValues(GuildID, ChannelID)
//	}
//
//	func ErroredChannelCalled(GuildID string, ChannelID string) prometheus.Counter {
//		return erroredChannels.WithLabelValues(GuildID, ChannelID)
//	}
func ChannelsPerGuild(GuildID string, ChannelID string) prometheus.Gauge {
	return channelsPerGuilds.WithLabelValues(GuildID, ChannelID)
}

func addMapGuildChannelValue(a_map map[string]map[string]int, guild string, channel_id string, value int) {
	guild_map, ok := a_map[guild]

	if !ok {
		guild_map = make(map[string]int)
		a_map[guild] = guild_map
	}

	guild_map[channel_id] += value
}

func Update(dg *discorder.Discorder, channels configurator.ConfiguratorChannel) {
	channelIDs, _ := channels.List()
	var channels_count_by_guild map[string]map[string]int = make(map[string]map[string]int)
	for _, channel := range channelIDs {
		_, err := dg.GetLatestMessages(channel)

		if ch, err := dg.GetDiscordSession().Channel(string(channel)); err == nil {
			if guild, err := dg.GetDiscordSession().Guild(ch.GuildID); err == nil {
				addMapGuildChannelValue(channels_count_by_guild, guild.Name, string(channel), 1)
			}
		}
		if err != nil {
			addMapGuildChannelValue(channels_count_by_guild, "unknown", string(channel), 1)
		}
	}

	DeleteAllGaugeVecValues(channelsPerGuilds)
	for guild_name, channels := range channels_count_by_guild {
		for channel_id, count := range channels {
			ChannelsPerGuild(guild_name, channel_id).Set(float64(count))
		}
	}
}

func Prometheuser() {
	dg := discorder.NewClient()
	channels := configurator.NewConfiguratorChannel(configurator.NewConfigurator(settings.Dbpath))
	go func() {
		// prometheus initializer
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8000", nil)
	}()

	for {
		Update(dg, channels)
		time.Sleep(time.Minute)
	}
}
