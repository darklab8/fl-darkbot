package viewer_msg

import (
	"github.com/darklab8/darklab_goutils/goutils/utils_logus"
)

func LogusMsg(value *Msg) utils_logus.SlogParam {
	return func(c *utils_logus.SlogGroup) {
		c.Params["msg_message_id"] = string(value.messageID)
		c.Params["msg_view_id"] = string(value.viewID)
		c.Params["msg_view_enumerated_id"] = string(value.viewEnumeratedID)
	}
}
