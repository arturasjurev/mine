package main

import (
	"log"
	"mine/factory"
	"mine/factory/grinder"
)

func main() {

	/*
		factory := &factory.Factory{
			Grinder: grinder.NewMemGrinder(),
		}
	*/

	m := factory.Mineral{
		Hardness: 1,
	}

	g := grinder.NewMemGrinder()

	if err := g.SetPower(1); err != nil {
		log.Fatal(err)
	}

	if err := g.Insert(m); err != nil {
		log.Fatal(err)
	}

	if err := g.Perform(); err != nil {
		log.Fatal(err)
	}

}
