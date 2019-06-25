package freezer

import (
	"context"
	"fmt"
	"time"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/minerals"
)

type MemFreezer struct {
	Power    int
	inserted bool
	resource minerals.Mineral
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

func (g *MemFreezer) Insert(item minerals.Mineral) error {
	if g.inserted {
		return fmt.Errorf("freezer is not empty")
	}
	g.resource = item
	g.inserted = true
	return nil
}

func (g *MemFreezer) Takeout() (minerals.Mineral, error) {
	g.inserted = false
	return g.resource, nil
}

func (g *MemFreezer) Process() error {
	if !g.inserted {
		return fmt.Errorf("freezer is empty")
	}
	processTime := factory.CalculateProcessTime(g.resource.Hardness, g.Power)

	if processTime != 0 {
		done := time.Tick(processTime)
		<-done
	}

	g.resource.State = minerals.Solid
	g.resource.Fractures = 0
	return nil
}

func (g *MemFreezer) ProcessWithCtx(ctx context.Context) error {
	processTime := factory.CalculateProcessTime(g.resource.Hardness, g.Power)
	done := time.Tick(processTime)

	select {
	case <-done:
		g.resource.State = minerals.Solid
		return nil
	case <-ctx.Done():
		return nil
	}
}
