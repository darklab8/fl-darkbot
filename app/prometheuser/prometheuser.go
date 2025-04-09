package prometheuser

import (
	"net/http"
	"time"

	"github.com/darklab8/fl-darkbot/app/configurator"
	"github.com/darklab8/fl-darkbot/app/discorder"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/collectors/version"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ver "github.com/prometheus/common/version"
)

var (
	channelsPerGuilds *prometheus.GaugeVec

	listenerAllowedOperations *prometheus.CounterVec

	viewerOperations *prometheus.CounterVec
)

func init() {
	// Create non-global registry.

	channelsPerGuilds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "darkbot_guilds_channels_gauge",
		Help: "The total number of channels used per guild",
	}, []string{"guild_name", "channel_id"})

	listenerAllowedOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "darkbot_listener_allowed_requests_count",
		Help: "Requests incoming to listener from Discord. Contains error if not allowed",
	}, []string{"guild_name", "channel_id", "error"})
	viewerOperations = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "darkbot_viewer_requests_count",
		Help: "Requests sent by viewer to handle table renderings. Contain error if smth went wrong",
	}, []string{"guild_name", "channel_id", "error"})

	newreg := prometheus.NewRegistry()
	reg := prometheus.WrapRegistererWith(
		prometheus.Labels{
			"environment": settings.Env.Environment,
			"version_id":  ver.Version,
		}, newreg)
	reg.MustRegister(
		channelsPerGuilds,
		listenerAllowedOperations,
		viewerOperations,
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		version.NewCollector("darkbot"),
	)

	http.Handle(
		"/metrics", promhttp.HandlerFor(
			newreg,
			promhttp.HandlerOpts{}),
	)
}

var old_labels map[string]map[string]int

func ChannelsPerGuild(GuildID string, ChannelID string) prometheus.Gauge {
	return channelsPerGuilds.WithLabelValues(GuildID, ChannelID)
}
func ListenerIsAllowedOperations(GuildID string, ChannelID string, Error error) prometheus.Counter {
	var StrError string
	if Error != nil {
		StrError = Error.Error()
	}
	return listenerAllowedOperations.WithLabelValues(GuildID, ChannelID, StrError)
}
func ViewerOperations(GuildID string, ChannelID string, Error error) prometheus.Counter {
	var StrError string
	if Error != nil {
		StrError = Error.Error()
	}
	return viewerOperations.WithLabelValues(GuildID, ChannelID, StrError)
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
		logus.Log.CheckError(err, "prometheuser: unable to get msgs")

		if ch, err := dg.GetDiscordSession().Channel(string(channel)); err == nil {
			if guild, err := dg.GetDiscordSession().Guild(ch.GuildID); err == nil {
				addMapGuildChannelValue(channels_count_by_guild, guild.Name, string(channel), 1)
			}
		}
		if err != nil {
			addMapGuildChannelValue(channels_count_by_guild, "unknown", string(channel), 1)
		}
	}

	if old_labels != nil {
		for guild_name, channels := range old_labels {
			for channel_id, _ := range channels {
				channelsPerGuilds.DeleteLabelValues(guild_name, channel_id)
			}
		}
	}
	for guild_name, channels := range channels_count_by_guild {
		for channel_id, count := range channels {
			ChannelsPerGuild(guild_name, channel_id).Set(float64(count))
		}
	}

	old_labels = channels_count_by_guild
}

func Prometheuser(dg *discorder.Discorder) {
	channels := configurator.NewConfiguratorChannel(configurator.NewConfigurator(settings.Dbpath))

	go func() {
		for {
			Update(dg, channels)
			time.Sleep(time.Minute)
		}
	}()

	// http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe("0.0.0.0:8000", nil)
	logus.Log.CheckPanic(err, "unable to serve http server prometheuser")
}
