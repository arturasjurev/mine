package grinder

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sheirys/mine/factory"
)

// MemGrinder satisfies Equipment interface defined in factory/equipment.go
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
	g.empty = true
	return g.inserted, nil
}

func (g *MemGrinder) Perform() error {
	if g.empty {
		return fmt.Errorf("grinder is empty")
	}
	processTime := factory.CalculateProcessTime(g.inserted.Hardness, g.power)

	if processTime != 0 {
		done := time.Tick(processTime)
		<-done
	}

	g.inserted.State = factory.Dust
	return nil
}

func (g *MemGrinder) PerformWithCtx(ctx context.Context) error {
	if g.empty {
		return fmt.Errorf("grinder is empty")
	}
	processTime := factory.CalculateProcessTime(g.inserted.Hardness, g.power)

	log.Printf("calculated time %s\n", processTime)
	done := time.Tick(processTime)

	select {
	case <-done:
		g.inserted.State = factory.Dust
		return nil
	case <-ctx.Done():
		return nil
	}
	return nil
}
