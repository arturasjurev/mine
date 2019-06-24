package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sheirys/mine/manager"
	"github.com/sheirys/mine/manager/journal"
)

var kills = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

func main() {

	amqpAddr := flag.String("a", "amqp://guest:guest@localhost:5672/", "rabbitmq connection")
	bindAddr := flag.String("b", "0.0.0.0:8080", "http listen bind")
	dataFile := flag.String("d", "data.json", "path to data file")

	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	m := &manager.Manager{
		Journal: &journal.JournalFileService{
			File: *dataFile,
		},
		AMQPAddress: *amqpAddr,
		HTTPAddress: *bindAddr,
	}

	m.Init()

	m.Start()

	<-stop

	m.Stop()

}
