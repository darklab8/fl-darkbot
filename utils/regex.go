package utils

import (
	"darkbot/utils/logger"
	"regexp"
)

func InitRegexExpression(regex **regexp.Regexp, expression string) {
	var err error

	*regex, err = regexp.Compile(expression)
	logger.CheckPanic(err, "failed to init regex={%s} in ", expression, GetCurrentFile())
}
