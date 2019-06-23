package grinder

import (
	"context"
	"fmt"
	"time"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/minerals"
)

// MemGrinder satisfies Equipment interface defined in factory/equipment.go
type MemGrinder struct {
	Power    int
	inserted bool
	resource minerals.Mineral
}

func NewMemGrinder() *MemGrinder {
	return &MemGrinder{}
}

func (g *MemGrinder) Empty() bool {
	return !g.inserted
}

func (g *MemGrinder) SetPower(watts int) error {

	// omg omg negtive power detected. this is grinder
	// not some interstellar warp engine.
	if watts < 0 {
		return fmt.Errorf("nagetive power")
	}
	g.Power = watts
	return nil
}

func (g *MemGrinder) Insert(item minerals.Mineral) error {
	if g.inserted {
		return fmt.Errorf("grinder is not empty")
	}
	g.resource = item
	g.inserted = true
	return nil
}

func (g *MemGrinder) Takeout() (minerals.Mineral, error) {
	g.inserted = false
	return g.resource, nil
}

func (g *MemGrinder) Process() error {
	if !g.inserted {
		return fmt.Errorf("grinder is empty")
	}
	processTime := factory.CalculateProcessTime(g.resource.Hardness, g.Power)

	if processTime != 0 {
		done := time.Tick(processTime)
		<-done
	}

	g.resource.State = minerals.Fracture
	if g.resource.Fractures > 0 {
		g.resource.Fractures *= 2
	} else {
		g.resource.Fractures = 1
	}
	return nil
}

func (g *MemGrinder) ProcessWithCtx(ctx context.Context) error {
	if !g.inserted {
		return fmt.Errorf("grinder is empty")
	}
	processTime := factory.CalculateProcessTime(g.resource.Hardness, g.Power)
	done := time.Tick(processTime)

	select {
	case <-done:
		g.resource.State = minerals.Fracture
		return nil
	case <-ctx.Done():
		return nil
	}
	return nil
}
