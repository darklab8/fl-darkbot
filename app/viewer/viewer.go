package viewer

import (
	"darkbot/app/scrappy"
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/settings/utils"
	"darkbot/app/settings/worker"
	"darkbot/app/viewer/apis"
	"time"
)

type ViewerDelays struct {
	betweenChannels types.ViewerDelayBetweenChannels
	betweenLoops    types.ViewerLoopDelay
}

type Viewer struct {
	delays  ViewerDelays
	api     *apis.API
	workers *worker.TaskPoolPeristent[*TaskRefreshChannel]
}

func NewViewer(dbpath types.Dbpath, scrappy_storage *scrappy.ScrappyStorage) *Viewer {
	api := apis.NewAPI(dbpath, scrappy_storage)
	v := &Viewer{
		api: api,
		delays: ViewerDelays{
			betweenChannels: 10,
			betweenLoops:    settings.ViewerLoopDelay,
		},
	}

	v.workers = worker.NewTaskPoolPersistent[*TaskRefreshChannel](
		worker.WithAllowFailedTasks[*TaskRefreshChannel](),
		worker.WithDisableParallelism[*TaskRefreshChannel](false),
		worker.WithWorkersAmount[*TaskRefreshChannel](10),
	)

	return v
}

func (v *Viewer) Run() {
	logus.Info("Viewer is now running.")

	go func() {
		for {
			v.workers.AwaitSomeTask()
		}
	}()
	for {
		v.Update()
	}
}

func (v *Viewer) Update() {
	time_viewer_started := time.Now()
	logus.Info("Viewer.Update")

	channelIDs, _ := v.api.Channels.List()
	logus.Info("Viewer.Update.channelIDs=", logus.ChannelIDs(channelIDs))

	// For each channel
	allChannelsTime := utils.NewTimeMeasure("all channels")
	for _, channelID := range channelIDs {
		utils.TimeMeasure(func() {
			task := NewRefreshChannelTask(v.api, channelID, v.delays.betweenChannels)
			// task.RunTask(worker_types.WorkerID(0))
			v.workers.DelayTask(task)
		}, "one channel", logus.ChannelID(channelID))
	}
	allChannelsTime.Close()
	logus.Info("Viewer.Update Finished " + time.Since(time_viewer_started).String())
	time.Sleep(time.Duration(v.delays.betweenLoops) * time.Second)
}
