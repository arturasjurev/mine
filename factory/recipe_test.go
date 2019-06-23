package factory_test

import (
	"testing"

	"github.com/sheirys/mine/factory"
	"github.com/sheirys/mine/minerals"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRecipe(t *testing.T) {
	testTable := []struct {
		From, To       minerals.State
		ExpectedRecipe []factory.RecipeAction
		ExpectedErr    bool
	}{
		{
			From: minerals.Fracture,
			To:   minerals.Solid,
			ExpectedRecipe: []factory.RecipeAction{
				factory.ApplySmelting,
				factory.ApplyFreezing,
			},
			ExpectedErr: false,
		},
		{
			From: minerals.Liquid,
			To:   minerals.Fracture,
			ExpectedRecipe: []factory.RecipeAction{
				factory.ApplyFreezing,
				factory.ApplyGrinding,
			},
			ExpectedErr: false,
		},
		{
			From: minerals.Fracture,
			To:   minerals.Fracture,
			ExpectedRecipe: []factory.RecipeAction{
				factory.ApplyGrinding,
			},
			ExpectedErr: false,
		},

		// expected to fail because of unkwnon states.
		{
			From:           "",
			To:             minerals.Fracture,
			ExpectedRecipe: []factory.RecipeAction{},
			ExpectedErr:    true,
		},
		{
			From:           minerals.Fracture,
			To:             "",
			ExpectedRecipe: []factory.RecipeAction{},
			ExpectedErr:    true,
		},
	}

	for i, v := range testTable {
		recipe, err := factory.GenerateRecipe(v.From, v.To)
		assert.Equal(t, v.ExpectedErr, err != nil, "case %d failed: %s", i, err)
		assert.Equal(t, v.ExpectedRecipe, recipe, "case %d failed", i)
	}
}
