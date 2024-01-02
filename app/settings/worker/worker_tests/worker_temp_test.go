package worker_tests

import (
	"darkbot/app/settings/darkbot_logus"
	"darkbot/app/settings/worker"
	"darkbot/app/settings/worker/worker_logus"
	"darkbot/app/settings/worker/worker_types"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/darklab8/darklab_goutils/goutils/logus"
	"github.com/stretchr/testify/assert"
)

// ======================
// Test Example

type TaskTest struct {
	*worker.Task

	// any desired arbitary data
	result worker_types.TaskID
}

func NewTaskTest(id worker_types.TaskID) *TaskTest {
	return &TaskTest{Task: worker.NewTask(id)}
}

func (data *TaskTest) RunTask(worker_id worker_types.WorkerID) worker_types.TaskStatusCode {
	darkbot_logus.Log.Debug("task test started", worker_logus.WorkerID(worker_id), worker_logus.TaskID(data.GetID()))
	time.Sleep(time.Second * time.Duration(data.GetID()))
	data.result = data.GetID() * 1
	data.SetAsDone()
	darkbot_logus.Log.Debug("task test finished", worker_logus.WorkerID(worker_id), worker_logus.TaskID(data.GetID()))
	return worker.CodeSuccess
}

func TaskResult(value worker_types.TaskID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["task_result"] = strconv.Itoa(int(value))
	}
}

func TestWorkerTemp(t *testing.T) {
	taskPool := worker.NewTaskPool[*TaskTest](
		worker.WithAllowFailedTasks[*TaskTest](),
		worker.WithDisableParallelism[*TaskTest](false),
	)

	tasks := []*TaskTest{}
	for task_id := 1; task_id <= 3; task_id++ {
		tasks = append(tasks, NewTaskTest(worker_types.TaskID(task_id)))
	}

	taskPool.RunTemporalPool(tasks)

	done_count := 0
	failed_count := 0
	for task_number, task := range tasks {
		darkbot_logus.Log.Debug(fmt.Sprintf("task.Done=%t", task.IsDone()), worker_logus.TaskID(worker_types.TaskID(task_number)), TaskResult(task.result))
		if task.IsDone() {
			done_count += 1
		} else {
			failed_count += 1
		}
	}
	assert.GreaterOrEqual(t, done_count, 3, "expected finding done tasks")
	assert.LessOrEqual(t, failed_count, 3, "expected finding failed tasks because of time sleep greater than timeout")
}
