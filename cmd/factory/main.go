package main

import (
	"os"
	"os/signal"
	"syscall"

)

var kills = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL}

func main() {

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, kills...)

	/*
		factory := &factory.Factory{
			Grinder: grinder.NewMemGrinder(),
		}
	*/

	/*
		m := factory.Mineral{
			Hardness: 150,
		}

		g := grinder.NewMemGrinder()

		if err := g.SetPower(3000); err != nil {
			log.Fatal(err)
		}

		if err := g.Insert(m); err != nil {
			log.Fatal(err)
		}

		if err := g.Process(); err != nil {
			log.Fatal(err)
		}
	*/

	/*
	worker := &dispatcher.Worker{
		Address: "amqp://guest:guest@localhost:5672/",
	}

	worker.Listen()
	*/
	<-stop
}
