package factory_test

import (
	"testing"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/factory/equipment/freezer"
	"github.com/sheirys/mine/factory/equipment/grinder"
	"github.com/sheirys/mine/factory/equipment/smelter"
	"github.com/stretchr/testify/assert"
)

func TestFactoryProcess(t *testing.T) {
	testTable := []struct {
		Resource        factory.Mineral
		From            factory.MineralState
		To              factory.MineralState
		ExpectedMineral factory.Mineral
		ExpectedErr     bool
	}{
		{
			Resource: factory.Mineral{
				Name:         "iron",
				State:        factory.Fracture,
				MeltingPoint: 2000,
				Hardness:     1000,
				Fractures:    2,
			},
			From: factory.Fracture,
			To:   factory.Solid,
			ExpectedMineral: factory.Mineral{
				Name:         "iron",
				State:        factory.Solid,
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
