package viewer_msg

import (
	"log/slog"

	"github.com/darklab8/go-utils/typelog"
)

func LogusMsg(value *Msg) typelog.LogType {
	return func(c *typelog.LogAtrs) {
		c.Append(
			slog.String("msg_message_id", string(value.messageID)),
			slog.String("msg_view_id", string(value.viewID)),
			slog.String("msg_view_enumerated_id", string(value.viewEnumeratedID)),
		)
	}
}
