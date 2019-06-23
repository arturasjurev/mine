package factory_test

import (
	"testing"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/factory/equipment/freezer"
	"github.com/sheirys/mine/factory/equipment/grinder"
	"github.com/sheirys/mine/factory/equipment/smelter"
	"github.com/sheirys/mine/minerals"
	"github.com/stretchr/testify/assert"
)

func TestFactoryProcess(t *testing.T) {
	testTable := []struct {
		Resource        minerals.Mineral
		From            minerals.State
		To              minerals.State
		ExpectedMineral minerals.Mineral
		ExpectedErr     bool
	}{
		{
			Resource: minerals.Mineral{
				Name:         "iron",
				State:        minerals.Fracture,
				MeltingPoint: 2000,
				Hardness:     1000,
				Fractures:    2,
			},
			From: minerals.Fracture,
			To:   minerals.Solid,
			ExpectedMineral: minerals.Mineral{
				Name:         "iron",
				State:        minerals.Solid,
				MeltingPoint: 2000,
				Hardness:     1000,
				Fractures:    0,
			},
			ExpectedErr: false,
		},
	}

	f := factory.Factory{
		Grinder: &grinder.MemGrinder{},
		Smelter: &smelter.MemSmelter{},
		Freezer: &freezer.MemFreezer{},
	}

	for i, v := range testTable {
		f.From = v.From
		f.To = v.To
		f.Resource = v.Resource

		err := f.Process()
		assert.Equal(t, v.ExpectedErr, err != nil, "case %d: %s", i, err)
		assert.Equal(t, v.ExpectedMineral, f.Resource)
	}

}
