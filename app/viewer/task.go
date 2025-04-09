package viewer

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/darklab8/fl-darkbot/app/prometheuser"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"
	"github.com/darklab8/go-typelog/typelog"

	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/worker"
	"github.com/darklab8/go-utils/utils/worker/worker_types"
)

type TaskRefreshChannel struct {
	*worker.Task

	// any desired arbitary data
	api                  *apis.API
	channelID            types.DiscordChannelID
	delayBetweenChannels types.ViewerDelayBetweenChannels
}

func NewRefreshChannelTask(
	api *apis.API,
	channelID types.DiscordChannelID,
	delayBetweenChannels types.ViewerDelayBetweenChannels,
) *TaskRefreshChannel {
	task_id_gen += 1
	return &TaskRefreshChannel{
		Task:                 worker.NewTask(worker_types.TaskID(task_id_gen)),
		api:                  api,
		channelID:            channelID,
		delayBetweenChannels: delayBetweenChannels,
	}
}

var task_id_gen int = 0

var guildAntiRateLimitMutexes map[string]*sync.Mutex

func init() {
	guildAntiRateLimitMutexes = make(map[string]*sync.Mutex)
}

func GetMutex(MutexKey string) *sync.Mutex {
	value, ok := guildAntiRateLimitMutexes[MutexKey]

	if ok {
		return value
	}

	new_mutex := &sync.Mutex{}
	guildAntiRateLimitMutexes[MutexKey] = new_mutex
	return new_mutex
}

func TryChannelAutoRemoval(api *apis.API, err error, channel_id types.DiscordChannelID) {
	logus_ch := logus.Log.WithFields(logus.ChannelID(channel_id), typelog.OptError(err))
	logus_ch.Info("checking channel for auto removing conditions")
	if strings.Contains(err.Error(), "Unknown Channel") {
		logus_ch.Info("channel matched condition to be removed because Unknown channel")
		logus_ch.CheckWarn(api.Channels.Remove(channel_id), "failed to delete channel")
		return
	}
	logus_ch.Warn("failed channel autoremoval")
}

func (v *TaskRefreshChannel) RunTask(worker_id worker_types.WorkerID) error {
	logus_ch := logus.Log.WithFields(logus.ChannelID(v.channelID))
	channel_info, err := v.api.Discorder.GetDiscordSession().Channel(string(v.channelID))

	var guild_name string = "unknown"
	if guild, _ := v.api.Discorder.GetDiscordSession().Guild(channel_info.GuildID); guild != nil {
		guild_name = guild.Name
	}

	if logus_ch.CheckError(err, "unable to get channel info") {
		prometheuser.ViewerOperations(guild_name, string(v.channelID), errors.New("unable to get channel info")).Inc()
		return err
	}

	MutexKey := channel_info.GuildID
	GuildMutex := GetMutex(MutexKey)
	GuildMutex.Lock()
	defer GuildMutex.Unlock()

	time_run_task_started := time.Now()
	time_new_channel := timeit.NewTimerL("new_channel", logus.ChannelID(v.channelID))
	channel := NewChannelView(v.api, v.channelID)

	time_new_channel.Close()

	time_render := timeit.NewTimerL("channel.Render", logus.ChannelID(v.channelID))
	channel.RenderViews()
	time_render.Close()

	time_discover := timeit.NewTimerL("channel.Discover", logus.ChannelID(v.channelID))
	err = channel.Discover()
	time_discover.Close()

	if logus_ch.CheckWarn(err, "unable to grab Discord msgs") {
		prometheuser.ViewerOperations(guild_name, string(v.channelID), errors.New("unable to grab Discord msgs")).Inc()
		return err
	}

	time_send := timeit.NewTimerL("channel.Send", logus.ChannelID(v.channelID))
	channel.Send()
	time_send.Close()

	time_delete_old := timeit.NewTimerL("channel.DeleteOld", logus.ChannelID(v.channelID))
	channel.DeleteOld()
	time_delete_old.Close()

	logus_ch.Info("RunTask finished",
		typelog.Int("task_id", int(v.Task.GetID())),
		typelog.String("elapsed", time.Since(time_run_task_started).String()),
		typelog.String("started_at", time_run_task_started.String()),
		typelog.String("finished_at", time.Now().String()),
	)

	prometheuser.ViewerOperations(guild_name, string(v.channelID), nil).Inc()

	// Important for Mutex above! Prevents Guild level rate limits. looks like 5 msg edits per 5 second at one server is good
	time.Sleep(time.Duration(v.delayBetweenChannels) * time.Second)
	return nil
}
