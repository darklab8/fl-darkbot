package worker

import (
	"github.com/darklab8/go-utils/utils/worker/worker_types"
)

type TaskPoolPeristent struct {
	*TaskPool

	task_channel chan ITask
}

func NewTaskPoolPersistent(opts ...TaskPoolOption) *TaskPoolPeristent {
	j := &TaskPoolPeristent{TaskPool: NewTaskPool(opts...)}

	j.task_channel = make(chan ITask)

	// This starts up N workers, initially blocked because there are no tasks yet.
	for worker_id := 1; worker_id <= j.numWorkers; worker_id++ {
		go j.launchWorker(
			worker_types.WorkerID(worker_id),
			j.task_channel,
		)
	}

	return j
}

func (j *TaskPoolPeristent) DelayTask(task ITask) {
	j.task_channel <- task
}
