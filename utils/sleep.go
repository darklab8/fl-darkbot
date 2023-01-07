package utils

import (
	"fmt"
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
	fmt.Println("awaiting smth forever")
	mu := sync.Mutex{}
	mu.Lock()
	mu.Lock()
}
