package freezer_test

import (
	"testing"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/factory/equipment/freezer"
	"github.com/sheirys/mine/minerals"
	"github.com/stretchr/testify/assert"
)

func TestMemFreezerEmpty(t *testing.T) {
	var (
		m   minerals.Mineral
		g   factory.Equipment
		err error
	)

	m = minerals.Mineral{}
	g = freezer.NewMemFreezer()

	// first, freezer should be empty.
	assert.Equal(t, true, g.Empty())

	// we should successfuly insert mineral into empty freezer.
	err = g.Insert(m)
	assert.NoError(t, err)

	// after insertion, freezer should not be empty.
	assert.Equal(t, false, g.Empty())

	// if freezer is not empty, we can takeout whatever is there.
	_, err = g.Takeout()
	assert.NoError(t, err)

	// after takeout, freezer should become empty again.
	assert.Equal(t, true, g.Empty())
}

// TestMemFreezerProcess check if mineral changes state after performing freeze on it.
func TestMemFreezerProcess(t *testing.T) {
	var (
		m, p minerals.Mineral
		g    factory.Equipment
		err  error
	)

	m = minerals.Mineral{
		Name:  "iron",
		State: minerals.Liquid,
	}
	g = freezer.NewMemFreezer()

	// process should fail on empty freezer.
	err = g.Process()
	assert.Error(t, err)

	// we should successfuly insert mineral into empty freezer.
	err = g.Insert(m)
	assert.NoError(t, err)

	// perform should not throw error, if it is not empty.
	err = g.Process()
	assert.NoError(t, err)

	// after performing, we should take out solid minerals.
	p, err = g.Takeout()
	assert.NoError(t, err)

	// after freezer, mineral state should be solid.
	assert.Equal(t, minerals.State(minerals.Solid), p.State)
}

func TestMemFreezerPower(t *testing.T) {
	g := freezer.NewMemFreezer()

	err := g.SetPower(-1)
	assert.Error(t, err)

	err = g.SetPower(0)
	assert.NoError(t, err)

	err = g.SetPower(100)
	assert.NoError(t, err)
}
