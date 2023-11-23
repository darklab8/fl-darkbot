package worker_logus

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/utils"
	"darkbot/app/viewer/worker/worker_types"
	"fmt"
	"strconv"
	"strings"
)

func WorkerID(value worker_types.WorkerID) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["worker_id"] = strconv.Itoa(int(value))
	}
}

func JobNumber(value int) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["job_number"] = strconv.Itoa(value)
	}
}

func StatusCodes(status_codes []worker_types.JobStatusCode) logus.SlogParam {
	str_status_codes := utils.CompL(status_codes, func(x worker_types.JobStatusCode) string { return fmt.Sprintf("%d", x) })
	return func(c *logus.SlogGroup) {

		c.Params["status_codes"] = strings.Join(str_status_codes, ",")
	}
}

func JobResult(value int) logus.SlogParam {
	return func(c *logus.SlogGroup) {
		c.Params["job_result"] = strconv.Itoa(value)
	}
}
