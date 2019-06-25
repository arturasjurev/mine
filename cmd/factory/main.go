package main

import (
	"flag"
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

	amqpAddr := flag.String("a", "amqp://guest:guest@localhost:5672/", "rabbitmq connection")
	grinderPower := flag.Int("grinder", 500, "grinder power")
	smelterPower := flag.Int("smelter", 500, "smelter power")
	freezerPower := flag.Int("freezer", 500, "freezer power")

	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	f := &factory.Factory{
		Grinder: &grinder.MemGrinder{
			Power: *grinderPower,
		},
		Smelter: &smelter.MemSmelter{
			Power: *smelterPower,
		},
		Freezer: &freezer.MemFreezer{
			Power: *freezerPower,
		},
		AMQPAddress: *amqpAddr,
	}

	f.Init()

	f.Start()

	<-stop

	f.Stop()
}
