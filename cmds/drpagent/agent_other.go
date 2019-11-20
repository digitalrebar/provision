// +build !windows

package main

import (
	"log"
	"os"
	"syscall"
)

func serve(options opts) {
	ch := make(chan os.Signal)
	if err := createAgent(options, os.Stderr); err != nil {
		log.Fatalf("%v", err)
	}
	go func() {
		for {
			c := <-ch
			switch c {
			case syscall.SIGTERM, syscall.SIGINT:
				if err := stopAgent(); err != nil {
					log.Fatalf("error stopping agent: %v", err)
				}
				return
			}
		}
	}()
	if err := machineAgent.Run(); err != nil {
		log.Fatalf("Error running agent: %v", err)
	}
}
