package worker

import (
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/settings/worker/worker_logus"
	"darkbot/app/settings/worker/worker_types"
	"time"
)

type TaskPoolPeristent[taskT ITask] struct {
	*TaskPool[taskT]

	task_channel   chan taskT
	result_channel chan worker_types.TaskStatusCode
}

func NewTaskPoolPersistent[taskT ITask](opts ...TaskPoolOption[taskT]) *TaskPoolPeristent[taskT] {
	j := &TaskPoolPeristent[taskT]{TaskPool: NewTaskPool[taskT](opts...)}

	j.task_channel = make(chan taskT)
	j.result_channel = make(chan worker_types.TaskStatusCode)

	// This starts up N workers, initially blocked because there are no tasks yet.
	for worker_id := 1; worker_id <= j.numWorkers; worker_id++ {
		go j.launchWorker(
			worker_types.WorkerID(worker_id),
			j.task_channel,
			j.result_channel,
		)
	}
	return j
}

func (j *TaskPoolPeristent[taskT]) DelayTask(task taskT) {
	j.task_channel <- task
}

func (j *TaskPoolPeristent[taskT]) AwaitSomeTask() {
	select {
	case status_code := <-j.result_channel:
		darkbot_logus.Log.Debug("finished some task succesfully", worker_logus.StatusCode(status_code))
	case <-time.After(time.Duration(j.taskTimeout) * time.Second):
		// non zero exit by timeout
		darkbot_logus.Log.Error("finished tasks with", worker_logus.StatusCode(CodeTimeout)) // TODO add worker_logus.TaskNumber(worker_types.TaskID(task_number)
	}
}

func (j *TaskPoolPeristent[taskT]) AwaitTimeouts() {
	for {
		j.AwaitSomeTask()
	}
}
