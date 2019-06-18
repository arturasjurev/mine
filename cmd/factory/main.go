package main

import (
	"log"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/factory/grinder"
)

func main() {

	/*
		factory := &factory.Factory{
			Grinder: grinder.NewMemGrinder(),
		}
	*/

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

}
