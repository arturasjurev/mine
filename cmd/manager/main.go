package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sheirys/mine/manager"
	"github.com/sheirys/mine/manager/journal"
)

var kills = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	m := &manager.Manager{
		Journal: &journal.JournalFileService{
			File: "data.json",
		},
		AMQPAddress: "amqp://guest:guest@localhost:5672/",
		HTTPAddress: "0.0.0.0:8833",
	}

	m.Init()

	m.Start()

	<-stop

	m.Stop()

}
