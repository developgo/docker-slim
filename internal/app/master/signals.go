package app

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

var appContinueChan = make(chan struct{})
var appDoneChan = make(chan struct{})

var signals = []os.Signal{
	syscall.SIGUSR1,
}

func initSignalHandlers() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, signals...)
	log.Debugf("docker-slim: listening for signals - %+v", signals)
	go func() {
		for {
			select {
			case <-appDoneChan:
				return
			case sig := <-sigChan:
				switch sig {
				case syscall.SIGUSR1:
					log.Debug("docker-slim: continue signal")
					appContinueChan <- struct{}{}
				default:
					log.Debugf("docker-slim: other signal (%v)...", sig)
				}
			}
		}
	}()
}
