package proc

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/windrivder/gopkg/logx"
)

func AwaitSignal() {
	c := make(chan os.Signal, 1)
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-c:
		logx.WithFields(logx.Fields{"signal": s.String()}).Info("receive a signal")
	}
}
