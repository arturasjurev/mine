package grinder_test

import (
	"testing"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/factory/equipment/grinder"
	"github.com/sheirys/mine/minerals"
	"github.com/stretchr/testify/assert"
)

// TestMemGrinderInsert check if grinder changes status from empty to not
// empty, after inserting minerals into it.
func TestMemGrinderEmpty(t *testing.T) {
	var (
		m   minerals.Mineral
		g   factory.Equipment
		err error
	)

	m = minerals.Mineral{}
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

func TestMemGrinderPower(t *testing.T) {
	g := grinder.NewMemGrinder()

	err := g.SetPower(-1)
	assert.Error(t, err)

	err = g.SetPower(0)
	assert.NoError(t, err)

	err = g.SetPower(100)
	assert.NoError(t, err)
}

func TestMemGrinderFractures(t *testing.T) {
	m := minerals.Mineral{
		Name:      "iron",
		State:     minerals.Fracture,
		Fractures: 4,
	}
	g := grinder.NewMemGrinder()

	err := g.Insert(m)
	assert.NoError(t, err)

	err = g.Process()
	assert.NoError(t, err)

	p, err := g.Takeout()
	assert.NoError(t, err)

	// after grinder, mineral state should be fracured.
	assert.Equal(t, minerals.State(minerals.Fracture), p.State)
	assert.Equal(t, 8, p.Fractures)
}
