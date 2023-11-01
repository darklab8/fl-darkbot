package utils

import (
	"darkbot/app/settings/logus"
	"darkbot/app/settings/types"
	"regexp"
)

func InitRegexExpression(regex **regexp.Regexp, expression types.RegExp) {
	var err error

	*regex, err = regexp.Compile(string(expression))
	logus.CheckFatal(err, "failed to init regex", logus.Regex(expression), logus.FilePath(GetCurrentFile()))
}
