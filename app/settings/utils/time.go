package utils

import (
	"darkbot/app/settings/logus"
	"fmt"
	"time"
)

type timeMeasurer struct {
	msg          string
	ops          []logus.SlogParam
	time_started time.Time
}

func NewTimeMeasure(msg string, ops ...logus.SlogParam) *timeMeasurer {
	return &timeMeasurer{
		msg:          msg,
		ops:          ops,
		time_started: time.Now(),
	}
}

func (t *timeMeasurer) Close() {
	logus.Debug(fmt.Sprintf("time_measure %v | %s", time.Since(t.time_started), t.msg), t.ops...)
}

func TimeMeasure(callback func(), msg string, ops ...logus.SlogParam) {
	time_started := NewTimeMeasure(msg, ops...)
	defer time_started.Close()
	callback()
}
