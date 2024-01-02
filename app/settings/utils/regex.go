package utils

import (
	"darkbot/app/settings/darkbot_logus"
	"regexp"

	"github.com/darklab8/darklab_goutils/goutils/utils_logus"
	"github.com/darklab8/darklab_goutils/goutils/utils_logus/logus_types"
)

func InitRegexExpression(regex **regexp.Regexp, expression logus_types.RegExp) {
	var err error

	*regex, err = regexp.Compile(string(expression))
	darkbot_logus.Log.CheckFatal(err, "failed to init regex", utils_logus.Regex(expression), utils_logus.FilePath(GetCurrentFile()))
}
