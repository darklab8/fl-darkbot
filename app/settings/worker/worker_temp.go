package worker

import (
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/settings/worker/worker_logus"
	"darkbot/app/settings/worker/worker_types"
	"time"
)

// ====================

type ITask interface {
	RunTask(worker_id worker_types.WorkerID) worker_types.TaskStatusCode
	IsDone() bool
}

type Task struct {
	id   worker_types.TaskID
	done bool
}

func (t *Task) GetID() worker_types.TaskID { return t.id }
func (t *Task) IsDone() bool               { return t.done }
func (t *Task) SetAsDone()                 { t.done = true }

func NewTask(id worker_types.TaskID) *Task {
	return &Task{id: id}
}

const (
	CodeSuccess worker_types.TaskStatusCode = 0
	CodeTimeout worker_types.TaskStatusCode = 1
	CodeFailure worker_types.TaskStatusCode = 2
)

type TaskPool[taskT ITask] struct {
	taskTimeout worker_types.Seconds
	numWorkers  int

	allow_failed_tasks  bool
	disable_parallelism worker_types.DebugDisableParallelism
}

type TaskPoolOption[T ITask] func(r *TaskPool[T])

func WithAllowFailedTasks[T ITask]() TaskPoolOption[T] {
	return func(c *TaskPool[T]) {
		c.allow_failed_tasks = true
	}
}

func WithWorkersAmount[T ITask](value int) TaskPoolOption[T] {
	return func(c *TaskPool[T]) { c.numWorkers = value }
}

func WithTaskTimeout[T ITask](value worker_types.Seconds) TaskPoolOption[T] {
	return func(c *TaskPool[T]) { c.taskTimeout = value }
}

func WithDisableParallelism[T ITask](disable_parallelism worker_types.DebugDisableParallelism) TaskPoolOption[T] {
	return func(c *TaskPool[T]) { c.disable_parallelism = disable_parallelism }
}

func NewTaskPool[T ITask](opts ...TaskPoolOption[T]) *TaskPool[T] {
	j := &TaskPool[T]{
		numWorkers:  3,
		taskTimeout: 120,
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

func (j *TaskPool[taskT]) launchWorker(worker_id worker_types.WorkerID, tasks <-chan taskT, results chan<- worker_types.TaskStatusCode) {
	darkbot_logus.Log.Info("worker started", worker_logus.WorkerID(worker_id))
	for task := range tasks {
		results <- task.RunTask(worker_id)
	}
	darkbot_logus.Log.Info("worker finished", worker_logus.WorkerID(worker_id))
}

/// Temporal

func (j *TaskPool[taskT]) runTasksinTemporalWorkers(tasks []taskT) []worker_types.TaskStatusCode {
	numTasks := len(tasks)

	// In order to use our pool of workers we need to send them work and collect their results.
	// We make 2 channels for this.
	tasks_channel := make(chan taskT, numTasks)
	result_channel := make(chan worker_types.TaskStatusCode, numTasks)
	status_codes := []worker_types.TaskStatusCode{}

	// This starts up N workers, initially blocked because there are no tasks yet.
	for worker_id := 1; worker_id <= j.numWorkers; worker_id++ {
		go j.launchWorker(worker_types.WorkerID(worker_id), tasks_channel, result_channel)
	}

	// Here we send 5 tasks and
	for _, task := range tasks {
		tasks_channel <- task
	}

	// then close that channel to indicate that is all the work we have.
	close(tasks_channel)

	// Finally we collect all the results of the work.
	// This also ensures that the worker goroutines have finished.
	// An alternative way to wait for multiple goroutines is to use a WaitGroup.
	for task_number := range tasks {
		select {
		case res := <-result_channel:
			status_codes = append(status_codes, res)
		case <-time.After(time.Duration(j.taskTimeout) * time.Second):
			// non zero exit by timeout
			status_codes = append(status_codes, CodeTimeout)
			darkbot_logus.Log.Error("timeout for", worker_logus.TaskID(worker_types.TaskID(task_number)))
		}

	}
	return status_codes
}

func (taskPool *TaskPool[taskT]) RunTemporalPool(tasks []taskT) {
	/*
		Switcher executing tasks with smth resembling multithreaded pool
		or executing synrconously if debug is on
	*/

	if taskPool.disable_parallelism {
		for pseudo_worker_id, task := range tasks {
			task.RunTask(worker_types.WorkerID(pseudo_worker_id))
		}
	} else {
		status_codes := taskPool.runTasksinTemporalWorkers(tasks)
		darkbot_logus.Log.Debug("results", worker_logus.StatusCodes(status_codes))
	}

	for task_number, task := range tasks {
		task_id := worker_types.TaskID(task_number)
		if !task.IsDone() && !taskPool.allow_failed_tasks {
			darkbot_logus.Log.Error("task failed", worker_logus.TaskID(task_id))
		}
		darkbot_logus.Log.Debug("task succeed", worker_logus.TaskID(task_id))
	}
}
