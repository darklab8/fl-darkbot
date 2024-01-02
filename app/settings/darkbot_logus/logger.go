package darkbot_logus

import (
	"darkbot/app/settings/envs"

	"github.com/darklab8/darklab_goutils/goutils/utils_logus"
	"github.com/darklab8/darklab_goutils/goutils/utils_logus/logus_types"
)

var (
	Log *utils_logus.Logger
)

func init() {
	Log = utils_logus.NewLogger(
		logus_types.LogLevel(envs.LogLevel),
		logus_types.EnableJsonFormat(envs.LogTurnJSONLogging),
		logus_types.EnableFileShowing(envs.LogShowFileLocations))
}
