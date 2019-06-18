package smelter

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/sheirys/mine/factory"
)

type MemSmelter struct {
	empty    bool
	power    int
	inserted factory.Mineral
}

func NewMemSmelter() *MemSmelter {
	return &MemSmelter{
		empty: true,
	}
}

func (g *MemSmelter) Empty() bool {
	return g.empty
}

func (g *MemSmelter) SetPower(watts int) error {
	if watts < 0 {
		return fmt.Errorf("nagetive power")
	}
	g.power = watts
	return nil
}

func (g *MemSmelter) Insert(item factory.Mineral) error {
	if !g.empty {
		return fmt.Errorf("smelter is not empty")
	}
	g.inserted = item
	g.empty = false
	return nil
}

func (g *MemSmelter) Takeout() (factory.Mineral, error) {
	return g.inserted, nil
}

func (g *MemSmelter) Perform() error {
	if g.empty {
		return fmt.Errorf("smelter is empty")
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

func (g *MemSmelter) PerformWithCtx(context.Context) error {
	return nil
}
