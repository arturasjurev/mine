package freezer

import (
	"context"
	"fmt"
	"time"

	"github.com/sheirys/mine/factory"
)

type MemFreezer struct {
	Power    int
	inserted bool
	resource factory.Mineral
}

func NewMemFreezer() *MemFreezer {
	return &MemFreezer{}
}

func (g *MemFreezer) Empty() bool {
	return !g.inserted
}

func (g *MemFreezer) SetPower(watts int) error {
	if watts < 0 {
		return fmt.Errorf("nagetive power")
	}
	g.Power = watts
	return nil
}

func (g *MemFreezer) Insert(item factory.Mineral) error {
	if !g.inserted {
		return fmt.Errorf("smelter is not empty")
	}
	g.resource = item
	g.inserted = true
	return nil
}

func (g *MemFreezer) Takeout() (factory.Mineral, error) {
	return g.resource, nil
}

func (g *MemFreezer) Perform() error {
	if !g.inserted {
		return fmt.Errorf("smelter is empty")
	}
	processTime := factory.CalculateProcessTime(g.resource.Hardness, g.Power)

	if processTime != 0 {
		done := time.Tick(processTime)
		<-done
	}

	g.resource.State = factory.Solid
	g.resource.Fractures = 0
	return nil
}

func (g *MemFreezer) PerformWithCtx(ctx context.Context) error {
	processTime := factory.CalculateProcessTime(g.resource.Hardness, g.Power)
	done := time.Tick(processTime)

	select {
	case <-done:
		g.resource.State = factory.Solid
		return nil
	case <-ctx.Done():
		return nil
	}
	return nil
}
