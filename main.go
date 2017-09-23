package main

import (
	"github.com/justuswilhelm/iodine-runner/lib"

	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	s := lib.CreateSupervisor()
	s.ServeBackground()
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		s.Stop()
		os.Exit(0)
	}()

	for {
		runtime.Gosched()
	}
}
