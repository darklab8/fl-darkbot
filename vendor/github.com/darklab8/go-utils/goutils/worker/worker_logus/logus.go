package worker_logus

import (
	"log/slog"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/goutils/worker/worker_types"
)

var Log *typelog.Logger = typelog.NewLogger("worker")

func WorkerID(value worker_types.WorkerID) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(slog.Int("worker_id", int(value)))
	}
}

func TaskID(value worker_types.TaskID) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(slog.Int("task_id", int(value)))
	}
}

func StatusCode(value worker_types.TaskStatusCode) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(slog.Int("status_code", int(value)))
	}
}
