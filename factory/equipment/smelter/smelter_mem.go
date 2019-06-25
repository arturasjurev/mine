package smelter

import (
	"fmt"
	"time"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/minerals"
	"github.com/sirupsen/logrus"
)

type MemSmelter struct {
	Power    int
	inserted bool
	resource minerals.Mineral
}

func NewMemSmelter() *MemSmelter {
	return &MemSmelter{}
}

func (g *MemSmelter) Empty() bool {
	return !g.inserted
}

func (g *MemSmelter) SetPower(watts int) error {
	if watts < 0 {
		return fmt.Errorf("nagetive power")
	}
	g.Power = watts
	return nil
}

func (g *MemSmelter) GetPower() int {
	return g.Power
}

func (g *MemSmelter) Insert(item minerals.Mineral) error {
	if g.inserted {
		return fmt.Errorf("smelter is not empty")
	}
	g.resource = item
	g.inserted = true
	return nil
}

func (g *MemSmelter) Takeout() (minerals.Mineral, error) {
	g.inserted = false
	return g.resource, nil
}

func (g *MemSmelter) Process() error {
	if !g.inserted {
		return fmt.Errorf("smelter is empty")
	}
	processTime := factory.CalculateProcessTime(g.resource.MeltingPoint, g.Power)
	logrus.WithField("expected_time", processTime).Info("smelting")

	if processTime != 0 {
		done := time.Tick(processTime)
		<-done
	}

	g.resource.State = minerals.Liquid
	g.resource.Fractures = 0
	return nil
}
