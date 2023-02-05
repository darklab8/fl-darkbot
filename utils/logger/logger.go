package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type LogTag string

const (
	DEBUG LogTag = "DEBUG"
	INFO  LogTag = "INFO"
	WARN  LogTag = "DEBUG"
	FATAL LogTag = "ERROR"
	PANIC LogTag = "PANIC"
)

func GetCallingFile() string {
	GetTwiceParentFunctionLocation := 2
	_, filename, _, _ := runtime.Caller(GetTwiceParentFunctionLocation)
	filename = filepath.Base(filename)
	return fmt.Sprintf("f:%s ", filename)
}

func FormatTag(tag LogTag) string {
	return fmt.Sprintf("t:%s ", string(tag))
}

func print(v ...any) {
	if os.Getenv("LOGGING") == "false" {
		return
	}

	log.Print(fmt.Sprintln(v...))
}

func Debug(v ...any) {
	print(FormatTag(DEBUG), GetCallingFile(), fmt.Sprint(v...))
}

func Info(v ...any) {
	print(FormatTag(INFO), GetCallingFile(), fmt.Sprint(v...))
}

func Warn(v ...any) {
	print(FormatTag(WARN), GetCallingFile(), fmt.Sprint(v...))
}

func Fatal(v ...any) {
	print(FormatTag(FATAL), GetCallingFile(), fmt.Sprint(v...))
}

func Panic(v ...any) {
	print(FormatTag(PANIC), GetCallingFile(), fmt.Sprint(v...))
}

func CheckWarn(err error, v ...any) {
	if err == nil {
		return
	}

	Warn(err, fmt.Sprint(v...))
}

func CheckFatal(err error, v ...any) {
	if err == nil {
		return
	}

	Fatal(err, fmt.Sprint(v...))
}

func CheckPanic(err error, v ...any) {
	if err == nil {
		return
	}

	Panic(err, fmt.Sprint(v...))
}

func Debugf(format string, v ...any) {
	Debug(fmt.Sprintf(format, v...))
}

func Infof(format string, v ...any) {
	Info(fmt.Sprintf(format, v...))
}

func Warnf(format string, v ...any) {
	Warn(fmt.Sprintf(format, v...))
}

func Fatalf(format string, v ...any) {
	Fatal(fmt.Sprintf(format, v...))
}

func Panicf(format string, v ...any) {
	Panic(fmt.Sprintf(format, v...))
}
