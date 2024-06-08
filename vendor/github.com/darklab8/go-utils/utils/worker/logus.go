package worker

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/darklab8/go-typelog/typelog"
	"github.com/darklab8/go-utils/utils"
)

func LogusStatusCodes(tasks []ITask) typelog.LogType {
	str_status_codes := utils.CompL(tasks, func(x ITask) string { return fmt.Sprintf("%d", x.GetStatusCode()) })
	return func(c *typelog.LogAtrs) {
		c.Append(slog.String("status_codes", strings.Join(str_status_codes, ",")))
	}
}
