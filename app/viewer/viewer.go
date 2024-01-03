package viewer

import (
	"darkbot/app/scrappy"
	"darkbot/app/settings"
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"darkbot/app/viewer/apis"
	"time"

	"github.com/darklab8/darklab_goutils/goutils/worker"

	"github.com/darklab8/darklab_goutils/goutils/utils"
)

type ViewerDelays struct {
	betweenChannels types.ViewerDelayBetweenChannels
	betweenLoops    types.ViewerLoopDelay
}

type Viewer struct {
	delays  ViewerDelays
	api     *apis.API
	workers *worker.TaskPoolPeristent
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

	v.workers = worker.NewTaskPoolPersistent(
		worker.WithAllowFailedTasks(),
		worker.WithDisableParallelism(false),
		worker.WithWorkersAmount(10),
	)

	return v
}

func (v *Viewer) Run() {
	logus.Log.Info("Viewer is now running.")

	for {
		v.Update()
	}
}

func (v *Viewer) Update() {
	time_viewer_started := time.Now()
	logus.Log.Info("Viewer.Update")

	channelIDs, _ := v.api.Channels.List()
	logus.Log.Info("Viewer.Update.channelIDs=", logus.ChannelIDs(channelIDs))

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
	logus.Log.Info("Viewer.Update Finished " + time.Since(time_viewer_started).String())
	time.Sleep(time.Duration(v.delays.betweenLoops) * time.Second)
}
