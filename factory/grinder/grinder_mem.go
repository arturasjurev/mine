package grinder

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"mine/factory"
)

type MemGrinder struct {
	mtx      *sync.Mutex
	empty    bool
	power    int
	inserted factory.Mineral
}

func NewMemGrinder() *MemGrinder {
	return &MemGrinder{
		mtx:   &sync.Mutex{},
		empty: true,
	}
}

func (g *MemGrinder) Empty() bool {
	return g.empty
}

func (g *MemGrinder) SetPower(watts int) error {

	// omg omg negtive power detected. this is grinder
	// not some interstellar warp engine.
	if watts < 0 {
		return fmt.Errorf("nagetive power")
	}
	g.power = watts
	return nil
}

func (g *MemGrinder) Insert(item factory.Mineral) error {
	if !g.empty {
		return fmt.Errorf("grinder is not empty")
	}
	g.inserted = item
	g.empty = false
	return nil
}

func (g *MemGrinder) Takeout() (factory.Mineral, error) {
	return g.inserted, nil
}

func (g *MemGrinder) Perform() error {
	if g.empty {
		return fmt.Errorf("grinder is empty")
	}
	if g.power > 0 {
		full := g.inserted.Hardness / g.power
		partial := g.inserted.Hardness % g.power
		processTime := time.Duration(full) * time.Second
		processTime += time.Duration(partial) * time.Millisecond

		log.Printf("calculated time %s\n", processTime)
		done := time.Tick(processTime)
		<-done
	}
	return nil
}

func (g *MemGrinder) PerformWithCtx(context.Context) error {
	return nil
}
