package factory_test

import (
	"testing"

	"github.com/sheirys/mine/factory"
	"github.com/stretchr/testify/assert"
)

func TestGenerateRecipe(t *testing.T) {
	testTable := []struct {
		From, To       factory.MineralState
		ExpectedRecipe []factory.RecipeAction
		ExpectedErr    bool
	}{
		{
			From: factory.Fracture,
			To:   factory.Solid,
			ExpectedRecipe: []factory.RecipeAction{
				factory.ApplySmelting,
				factory.ApplyFreezing,
			},
			ExpectedErr: false,
		},
		{
			From: factory.Liquid,
			To:   factory.Fracture,
			ExpectedRecipe: []factory.RecipeAction{
				factory.ApplyFreezing,
				factory.ApplyGrinding,
			},
			ExpectedErr: false,
		},
		{
			From: factory.Fracture,
			To:   factory.Fracture,
			ExpectedRecipe: []factory.RecipeAction{
				factory.ApplyGrinding,
			},
			ExpectedErr: false,
		},
	}

	for i, v := range testTable {
		recipe, err := factory.GenerateRecipe(v.From, v.To)
		assert.Equal(t, v.ExpectedErr, err != nil, "case %d failed: %s", i, err)
		assert.Equal(t, v.ExpectedRecipe, recipe, "case %d failed", i)
	}
}
