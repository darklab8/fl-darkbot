package darkbot_logus

import (
	"darkbot/app/settings/envs"

	"github.com/darklab8/darklab_goutils/goutils/logus"
	"github.com/darklab8/darklab_goutils/goutils/logus/logus_types"
)

var (
	Log *logus.Logger
)

func init() {
	Log = logus.NewLogger(
		logus_types.LogLevel(envs.LogLevel),
		logus_types.EnableJsonFormat(envs.LogTurnJSONLogging),
		logus_types.EnableFileShowing(envs.LogShowFileLocations))
}
