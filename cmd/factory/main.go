package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/factory/equipment/freezer"
	"github.com/sheirys/mine/factory/equipment/grinder"
	"github.com/sheirys/mine/factory/equipment/smelter"
)

var kills = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	f := &factory.Factory{
		Grinder: &grinder.MemGrinder{
			Power: 500,
		},
		Smelter: &smelter.MemSmelter{
			Power: 500,
		},
		Freezer: &freezer.MemFreezer{
			Power: 500,
		},
		AMQPAddress: "amqp://guest:guest@localhost:5672/",
	}

	f.Init()

	f.Start()

	<-stop

	f.Stop()
}
