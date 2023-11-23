package worker_tests

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/worker"
	"darkbot/app/settings/worker/worker_logus"
	"darkbot/app/settings/worker/worker_types"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPersistent(t *testing.T) {
	taskPool := worker.NewTaskPoolPersistent[*TaskTest](
		worker.WithAllowFailedTasks[*TaskTest](),
		worker.WithDisableParallelism[*TaskTest](false),
	)

	tasks := []*TaskTest{}
	for task_id := 1; task_id <= 3; task_id++ {
		tasks = append(tasks, NewTaskTest(worker_types.TaskID(2)))
	}

	for _, task := range tasks {
		taskPool.DelayTask(task)
	}

	// U can test that it works even without awaitings
	// Awaiting is during prod usage necessary if u are going to timeout tasks though
	for range tasks {
		taskPool.AwaitSomeTask()
	}

	done_count := 0
	failed_count := 0
	for task_number, task := range tasks {
		logus.Debug(fmt.Sprintf("task.Done=%t", task.IsDone()), worker_logus.TaskID(worker_types.TaskID(task_number)), TaskResult(task.result))
		if task.IsDone() {
			done_count += 1
		} else {
			failed_count += 1
		}
	}
	assert.GreaterOrEqual(t, done_count, 3, "expected finding done tasks")
	assert.LessOrEqual(t, failed_count, 3, "expected finding failed tasks because of time sleep greater than timeout")
}
