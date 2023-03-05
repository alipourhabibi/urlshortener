package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/alipourhabibi/urlshortener/cmd"
)

func main() {
	var state byte
	const (
		reconfigure byte = iota
		waitForSignal
	)
	signalCh := make(chan os.Signal, 1)
	for {
		switch state {
		case reconfigure:
			state = waitForSignal
			go func() {
				// Entrypoint for application
				cmd.Execute()
				signalCh <- syscall.SIGQUIT
			}()
		case waitForSignal:
			signal.Notify(signalCh,
				syscall.SIGHUP,
				syscall.SIGINT,
				syscall.SIGTERM,
				syscall.SIGQUIT)

			sig := <-signalCh
			log.Println("signal recieved: ", sig)
			switch sig {
			case syscall.SIGHUP:
				state = reconfigure
			case syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM:
				return
			}

		}
	}

}
