package worker_logus

import (
	"darkbot/app/settings/worker/worker_types"
	"fmt"
	"strconv"
	"strings"

	"github.com/darklab8/darklab_goutils/goutils/utils"

	"github.com/darklab8/darklab_goutils/goutils/logus"
)

func WorkerID(value worker_types.WorkerID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["worker_id"] = strconv.Itoa(int(value))
	}
}

func TaskID(value worker_types.TaskID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["task_id"] = strconv.Itoa(int(value))
	}
}

func StatusCodes(status_codes []worker_types.TaskStatusCode) logus.SlogParam {
	str_status_codes := utils.CompL(status_codes, func(x worker_types.TaskStatusCode) string { return fmt.Sprintf("%d", x) })
	return func(c *logus.SlogGroup) {

		c.Params["status_codes"] = strings.Join(str_status_codes, ",")
	}
}

func StatusCode(value worker_types.TaskStatusCode) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["status_code"] = strconv.Itoa(int(value))
	}
}
