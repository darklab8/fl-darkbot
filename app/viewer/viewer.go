package viewer

import (
	"time"

	"github.com/darklab8/fl-darkbot/app/scrappy"
	"github.com/darklab8/fl-darkbot/app/settings"
	"github.com/darklab8/fl-darkbot/app/settings/logus"
	"github.com/darklab8/fl-darkbot/app/settings/types"
	"github.com/darklab8/fl-darkbot/app/viewer/apis"

	"github.com/darklab8/go-utils/utils/timeit"
	"github.com/darklab8/go-utils/utils/worker"
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
	allChannelsTime := timeit.NewTimer("all channels")
	for _, channelID := range channelIDs {
		timeit.NewTimerF(func(m *timeit.Timer) {
			task := NewRefreshChannelTask(v.api, channelID, v.delays.betweenChannels)
			// task.RunTask(worker_types.WorkerID(0))
			v.workers.DelayTask(task)
		}, timeit.WithMsg("one channel"), timeit.WithLogs(logus.ChannelID(channelID)))
	}
	allChannelsTime.Close()
	logus.Log.Info("Viewer.Update Finished " + time.Since(time_viewer_started).String())
	time.Sleep(time.Duration(v.delays.betweenLoops) * time.Second)
}
