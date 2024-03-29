package viewer

import (
	"fmt"
	"sync"
	"time"

	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"

	"github.com/darklab8/go-utils/goutils/worker"
	"github.com/darklab8/go-utils/goutils/worker/worker_types"

	"github.com/darklab8/go-utils/goutils/utils"
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

func (v *TaskRefreshChannel) RunTask(worker_id worker_types.WorkerID) error {
	channel_info, err := v.api.Discorder.GetDiscordSession().Channel(string(v.channelID))
	if logus.Log.CheckError(err, "unable to get channel info", logus.ChannelID(v.channelID)) {
		return err
	}

	MutexKey := channel_info.GuildID
	GuildMutex := GetMutex(MutexKey)
	GuildMutex.Lock()
	defer GuildMutex.Unlock()

	time_run_task_started := time.Now()
	time_new_channel := utils.NewTimeMeasure("new_channel", logus.ChannelID(v.channelID))
	channel := NewChannelView(v.api, v.channelID)

	time_new_channel.Close()

	time_render := utils.NewTimeMeasure("channel.Render", logus.ChannelID(v.channelID))
	channel.RenderViews()
	time_render.Close()

	time_discover := utils.NewTimeMeasure("channel.Discover", logus.ChannelID(v.channelID))
	err = channel.Discover()
	time_discover.Close()

	if logus.Log.CheckWarn(err, "unable to grab Discord msgs", logus.ChannelID(v.channelID)) {
		return err
	}

	time_send := utils.NewTimeMeasure("channel.Send", logus.ChannelID(v.channelID))
	channel.Send()
	time_send.Close()

	time_delete_old := utils.NewTimeMeasure("channel.DeleteOld", logus.ChannelID(v.channelID))
	channel.DeleteOld()
	time_delete_old.Close()

	logus.Log.Info(fmt.Sprintf("RunTask finished, TaskID=%d, elapsed=%s, started_at=%s, finished_at=%s",
		v.Task.GetID(),
		time.Since(time_run_task_started).String(),
		time_run_task_started.String(),
		time.Now().String(),
	))

	// Important for Mutex above! Prevents Guild level rate limits. looks like 5 msg edits per 5 second at one server is good
	time.Sleep(time.Duration(v.delayBetweenChannels) * time.Second)
	return nil
}
