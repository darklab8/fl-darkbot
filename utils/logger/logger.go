package logger

import (
	"fmt"
	"log"
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

func Debug(v ...any) {
	log.Print(FormatTag(DEBUG), GetCallingFile(), fmt.Sprintln(v...))
}

func Info(v ...any) {
	log.Print(FormatTag(INFO), GetCallingFile(), fmt.Sprintln(v...))
}

func Warn(v ...any) {
	log.Print(FormatTag(WARN), GetCallingFile(), fmt.Sprintln(v...))
}

func Fatal(v ...any) {
	log.Fatal(FormatTag(FATAL), GetCallingFile(), fmt.Sprintln(v...))
}

func Panic(v ...any) {
	log.Panic(FormatTag(PANIC), GetCallingFile(), fmt.Sprintln(v...))
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
