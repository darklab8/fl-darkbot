package envs

import (
	"os"
	"strings"

	"github.com/darklab8/darklab_goutils/goutils/logus_core/logus_types"
)

var LogTurnJSONLogging bool
var LogShowFileLocations bool
var LogLevel logus_types.LogLevel

const (
	EnvTrue = "true"
)

func init() {
	LogTurnJSONLogging = strings.ToLower(os.Getenv("github.com/darklab/fl-darkbot_LOG_JSON")) == EnvTrue
	LogShowFileLocations = strings.ToLower(os.Getenv("github.com/darklab/fl-darkbot_LOG_SHOW_FILE_LOCATIONS")) == EnvTrue

	log_level_str, is_log_level_set := os.LookupEnv("github.com/darklab/fl-darkbot_LOG_LEVEL")
	if !is_log_level_set {
		log_level_str = "WARN"
	}
	LogLevel = logus_types.LogLevel(log_level_str)
}
