package factory_test

import (
	"mine/factory"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRecipe(t *testing.T) {
	testTable := []struct {
		From, To       factory.Mineral
		ExpectedRecipe []factory.RecipeAction
		ExpectedErr    bool
	}{
		{
			From: factory.Mineral{State: factory.Fracture},
			To:   factory.Mineral{State: factory.Solid},
			ExpectedRecipe: []factory.RecipeAction{
				factory.ApplyGrinding,
				factory.ApplySmelting,
				factory.ApplyFreezing,
			},
			ExpectedErr: false,
		},
	}

	for i, v := range testTable {
		recipe, err := factory.GenerateRecipe(v.From, v.To)
		assert.Equal(t, v.ExpectedErr, err != nil, "case %d failed", i)
		assert.Equal(t, v.ExpectedRecipe, recipe, "case %d failed", i)
	}
}
