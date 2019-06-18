package grinder_test

import (
	"testing"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/factory/equipment/grinder"
	"github.com/stretchr/testify/assert"
)

// TestMemGrinderInsert check if grinder changes status from empty to not
// empty, after inserting minerals into it.
func TestMemGrinderEmpty(t *testing.T) {
	var (
		m   factory.Mineral
		g   factory.Equipment
		err error
	)

	m = factory.Mineral{}
	g = grinder.NewMemGrinder()

	// first, grinder should be empty.
	assert.Equal(t, true, g.Empty())

	// we should successfuly insert mineral into empty grinder.
	err = g.Insert(m)
	assert.NoError(t, err)

	// after insertion, grinder should not be empty.
	assert.Equal(t, false, g.Empty())

	// if grinder is not empty, we can takeout whatever is there.
	_, err = g.Takeout()
	assert.NoError(t, err)

	// after takeout, grinder should become empty again.
	assert.Equal(t, true, g.Empty())
}

// TestMemGrinderProcess check if mineral changes state after performing grind on it.
func TestMemGrinderProcess(t *testing.T) {
	var (
		m, p factory.Mineral
		g    factory.Equipment
		err  error
	)

	m = factory.Mineral{
		Name:  "iron",
		State: factory.Fracture,
	}
	g = grinder.NewMemGrinder()

	// process should fail on empty grinder.
	err = g.Process()
	assert.Error(t, err)

	// we should successfuly insert mineral into empty grinder.
	err = g.Insert(m)
	assert.NoError(t, err)

	// perform should not throw error, if it is not empty.
	err = g.Process()
	assert.NoError(t, err)

	// after performing, we should take out grinded minerals.
	p, err = g.Takeout()
	assert.NoError(t, err)

	// after grinder, mineral state should be fracured.
	assert.Equal(t, factory.MineralState(factory.Fracture), p.State)
}
