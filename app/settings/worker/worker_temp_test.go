package worker

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/worker/worker_logus"
	"darkbot/app/settings/worker/worker_types"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ======================
// Test Example

type TaskTest struct {
	*Task

	// any desired arbitary data
	result worker_types.TaskID
}

func NewTaskTest(id worker_types.TaskID) *TaskTest {
	return &TaskTest{Task: NewTask(id)}
}

func (data *TaskTest) runTask(worker_id worker_types.WorkerID) worker_types.TaskStatusCode {
	logus.Debug("task test started", worker_logus.WorkerID(worker_id), worker_logus.TaskID(data.id))
	time.Sleep(time.Second * time.Duration(data.id))
	data.result = data.id * 1
	data.done = true
	logus.Debug("task test finished", worker_logus.WorkerID(worker_id), worker_logus.TaskID(data.id))
	return CodeSuccess
}

func TaskResult(value worker_types.TaskID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["task_result"] = strconv.Itoa(int(value))
	}
}

func TestWorkerTemp(t *testing.T) {
	taskPool := NewTaskPool[*TaskTest](
		WithAllowFailedTasks[*TaskTest](),
		WithDisableParallelism[*TaskTest](false),
	)

	tasks := []*TaskTest{}
	for task_id := 1; task_id <= 3; task_id++ {
		tasks = append(tasks, NewTaskTest(worker_types.TaskID(task_id)))
	}

	taskPool.RunTemporalPool(tasks)

	done_count := 0
	failed_count := 0
	for task_number, task := range tasks {
		logus.Debug(fmt.Sprintf("task.Done=%t", task.done), worker_logus.TaskID(worker_types.TaskID(task_number)), TaskResult(task.result))
		if task.done {
			done_count += 1
		} else {
			failed_count += 1
		}
	}
	assert.GreaterOrEqual(t, done_count, 3, "expected finding done tasks")
	assert.LessOrEqual(t, failed_count, 3, "expected finding failed tasks because of time sleep greater than timeout")
}
