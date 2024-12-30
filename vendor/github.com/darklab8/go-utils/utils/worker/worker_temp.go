package worker

import (
	"errors"
	"fmt"
	"time"

	"github.com/darklab8/go-typelog/examples/logus"
	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils/worker/worker_logus"
	"github.com/darklab8/go-utils/utils/worker/worker_types"
)

// ====================

type ITask interface {
	RunTask(worker_id worker_types.WorkerID) error
	IsDone() bool
	setError(error)
	GetStatusCode() worker_types.TaskStatusCode
	SetAsDone()
}

type Task struct {
	id   worker_types.TaskID
	err  error
	done bool
}

func (t *Task) GetID() worker_types.TaskID { return t.id }
func (t *Task) IsDone() bool               { return t.done }

func (t *Task) SetAsDone()         { t.done = true }
func (t *Task) setError(err error) { t.err = err }

func (t *Task) GetStatusCode() worker_types.TaskStatusCode {
	if t.err == nil && t.done {
		return CodeSuccess
	}
	return CodeFailure
}

func NewTask(id worker_types.TaskID) *Task {
	return &Task{id: id}
}

const (
	CodeSuccess worker_types.TaskStatusCode = 0
	CodeTimeout worker_types.TaskStatusCode = 1
	CodeFailure worker_types.TaskStatusCode = 2
)

type TaskPool struct {
	taskTimeout worker_types.Seconds
	numWorkers  int

	allow_failed_tasks  bool
	disable_parallelism worker_types.DebugDisableParallelism

	task_observers []func(task ITask)
}

type TaskPoolOption func(r *TaskPool)

func WithTaskObServer(task_observer func(task ITask)) TaskPoolOption {
	return func(c *TaskPool) {
		c.task_observers = append(c.task_observers, task_observer)
	}
}

func WithAllowFailedTasks() TaskPoolOption {
	return func(c *TaskPool) {
		c.allow_failed_tasks = true
	}
}

func WithWorkersAmount(value int) TaskPoolOption {
	return func(c *TaskPool) { c.numWorkers = value }
}

func WithTaskTimeout(value worker_types.Seconds) TaskPoolOption {
	return func(c *TaskPool) { c.taskTimeout = value }
}

func WithDisableParallelism(disable_parallelism worker_types.DebugDisableParallelism) TaskPoolOption {
	return func(c *TaskPool) { c.disable_parallelism = disable_parallelism }
}

func NewTaskPool(opts ...TaskPoolOption) *TaskPool {
	j := &TaskPool{
		numWorkers:  3,
		taskTimeout: 60 * 30,
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

func (j *TaskPool) launchWorker(worker_id worker_types.WorkerID, tasks <-chan ITask) {
	worker_logus.Log.Info("worker started", worker_logus.WorkerID(worker_id))
	for task := range tasks {

		task_err := make(chan error, 1)
		go func() {
			defer func() {
				if !j.allow_failed_tasks {
					return
				}
				if r := recover(); r != nil {
					logus.Log.Error("Recovered in doRunf", typelog.Any("panic", r))
					task_err <- errors.New(fmt.Sprintln("task paniced", r))
				}
			}()
			task_err <- task.RunTask(worker_id)
		}()

		select {
		// Finally we collect all the results of the work.
		// This also ensures that the worker goroutines have finished.
		// An alternative way to wait for multiple goroutines is to use a WaitGroup.
		case <-time.After(time.Duration(j.taskTimeout) * time.Second):
			task.setError(errors.New("timed out"))
			worker_logus.Log.Error("finished tasks with timeout", worker_logus.StatusCode(CodeTimeout))
		case task_err := <-task_err:
			worker_logus.Log.CheckError(task_err, "finished tasks with timeout", worker_logus.StatusCode(CodeFailure))
			task.setError(task_err)
		}

		task.SetAsDone()
		for _, task_observer := range j.task_observers {
			task_observer(task)
		}
	}
	worker_logus.Log.Info("worker finished", worker_logus.WorkerID(worker_id))
}

/// Temporal

func runTasksinTemporalWorkers(tasks []ITask, j *TaskPool) {
	numTasks := len(tasks)

	// In order to use our pool of workers we need to send them work and collect their results.
	// We make 2 channels for this.
	tasks_channel := make(chan ITask, numTasks)

	// This starts up N workers, initially blocked because there are no tasks yet.
	for worker_id := 1; worker_id <= j.numWorkers; worker_id++ {
		go j.launchWorker(worker_types.WorkerID(worker_id), tasks_channel)
	}

	// Here we send 5 tasks and
	for _, task := range tasks {
		tasks_channel <- task
	}

	// then close that channel to indicate that is all the work we have.
	close(tasks_channel)
}

func RunTasksInTempPool(tasks []ITask, opts ...TaskPoolOption) {
	numTasks := len(tasks)
	result_channel := make(chan ITask, numTasks)

	total_options := []TaskPoolOption{
		WithTaskObServer(func(task ITask) {
			result_channel <- task
		}),
	}

	total_options = append(total_options, opts...)
	taskPool := NewTaskPool(total_options...)
	finished_tasks := []ITask{}

	/*
		Switcher executing tasks with smth resembling multithreaded pool
		or executing synrconously if debug is on
	*/
	if taskPool.disable_parallelism {
		for pseudo_worker_id, task := range tasks {
			task.RunTask(worker_types.WorkerID(pseudo_worker_id))
			finished_tasks = append(finished_tasks, task)
		}
	} else {
		runTasksinTemporalWorkers(tasks, taskPool)
		worker_logus.Log.Debug("results", LogusStatusCodes(tasks))

		for _ = range tasks {
			finished_tasks = append(finished_tasks, <-result_channel)
		}
	}

	for task_number, task := range tasks {
		task_id := worker_types.TaskID(task_number)
		if !task.IsDone() && !taskPool.allow_failed_tasks {
			worker_logus.Log.Error("task failed", worker_logus.TaskID(task_id))
		}
		worker_logus.Log.Debug("task succeed", worker_logus.TaskID(task_id))
	}
}
