package utils

import (
	"darkbot/settings/utils/logger"
	"fmt"
	"regexp"
)

func InitRegexExpression(regex **regexp.Regexp, expression string) {
	var err error

	*regex, err = regexp.Compile(expression)
	logger.CheckPanic(err, fmt.Sprintf("failed to init regex={%s} in ", expression), GetCurrentFile())
}
