package logger

import (
	"testing"
)

type CustomException struct {
}

func (c CustomException) Error() string {
	return "CustomException"
}

func TestLogging(t *testing.T) {
	Info("abc 123")
	Info("abc dfgd")
	Warn("123", "abc")
	Debug("fgdfg", "sdf")
	Info(CustomException{}.Error(), "abc")
	Info(CustomException{}.Error(), "abc")
}
