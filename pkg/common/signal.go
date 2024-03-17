package common

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func SigWatch(end context.CancelFunc, wait time.Duration, lg Logger) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-signals
	end()

	<-time.After(wait)
	lg.Sync()
}
