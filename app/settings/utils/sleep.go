package utils

import (
	"darkbot/app/settings/logus"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func SleepAwaitCtrlC() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func SleepForever() {
	logus.Debug("awaiting smth forever in SleepForever")
	mu := sync.Mutex{}
	mu.Lock()
	mu.Lock()
}