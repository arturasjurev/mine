package smelter_test

import (
	"testing"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/factory/equipment/smelter"
	"github.com/sheirys/mine/minerals"
	"github.com/stretchr/testify/assert"
)

func TestMemSmelterEmpty(t *testing.T) {
	var (
		m   minerals.Mineral
		g   factory.Equipment
		err error
	)

	m = minerals.Mineral{}
	g = smelter.NewMemSmelter()

	// first, smelter should be empty.
	assert.Equal(t, true, g.Empty())

	// we should successfuly insert mineral into empty smelter.
	err = g.Insert(m)
	assert.NoError(t, err)

	// after insertion, smelter should not be empty.
	assert.Equal(t, false, g.Empty())

	// if smelter is not empty, we can takeout whatever is there.
	_, err = g.Takeout()
	assert.NoError(t, err)

	// after takeout, smelter should become empty again.
	assert.Equal(t, true, g.Empty())
}

// TestMemSmelterProcess check if mineral changes state after performing smelt on it.
func TestMemSmelterProcess(t *testing.T) {
	var (
		m, p minerals.Mineral
		g    factory.Equipment
		err  error
	)

	m = minerals.Mineral{
		Name:  "iron",
		State: minerals.Fracture,
	}
	g = smelter.NewMemSmelter()

	// process should fail on empty smelter.
	err = g.Process()
	assert.Error(t, err)

	// we should successfuly insert mineral into empty smelter.
	err = g.Insert(m)
	assert.NoError(t, err)

	// perform should not throw error, if it is not empty.
	err = g.Process()
	assert.NoError(t, err)

	// after performing, we should take out liquid minerals.
	p, err = g.Takeout()
	assert.NoError(t, err)

	// after smelter, mineral state should be liquid.
	assert.Equal(t, minerals.State(minerals.Liquid), p.State)
	assert.Equal(t, 0, p.Fractures)
}

func TestMemSmelterPower(t *testing.T) {
	g := smelter.NewMemSmelter()

	err := g.SetPower(-1)
	assert.Error(t, err)

	err = g.SetPower(0)
	assert.NoError(t, err)

	err = g.SetPower(100)
	assert.NoError(t, err)
}
