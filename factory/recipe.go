package factory

import (
	"fmt"

	"github.com/sheirys/mine/minerals"
)

// This file contains factory logic and formulas how mineral states should be
// changed.

type RecipeAction string

// Define possible recept actions.
const (
	ApplyGrinding RecipeAction = "apply_grinding"
	ApplySmelting RecipeAction = "apply_smelting"
	ApplyFreezing RecipeAction = "apply_freezing"
)

// GenerateRecipe will define action order what equipments should be used and
// what actions should be applied to current mineral to reach asked mineral state.
// E.g.: if we pass current = fracture, asked = solid, this function should return
// recipe: []RecipeAction{"apply_grinding", "apply_smelting", "apply_freezing"}
func GenerateRecipe(current, asked minerals.State) ([]RecipeAction, error) {

	recipe := []RecipeAction{}

	// check if current and asked states is know for us. If states are
	// unknown for us, we do not know how to reach those states.
	// FIXME: minetals is now lib
	if _, ok := minerals.StateTable[current]; !ok {
		return recipe, fmt.Errorf("current state '%s' is unknown", current)
	}
	if _, ok := minerals.StateTable[asked]; !ok {
		return recipe, fmt.Errorf("asked state '%s' in unknown", asked)
	}

	return chainActions(recipe, current, asked)
}

// chainActions is core logic how to chain recipe actions. This will return
// recipe in form []RecipeAction{"apply_grinding", "apply_smelting", "apply_freezing"}.
// This function is used inside GenerateRecipe function.
// SPOILER ALERT: RECURSION USED
func chainActions(recipe []RecipeAction, now, stop minerals.State) ([]RecipeAction, error) {

	// NOTICE: there is special case. Only on fractured state it is
	// possible to apply grinding again, because sometimes what we
	// want to achieve is to double its fractures. Because of this
	// we can process order like `order{from: fractured, to: fractured}`
	// in case we want to double its fractures.
	if now == stop && stop == minerals.Fracture && len(recipe) == 0 {
		recipe = append(recipe, ApplyGrinding)
		return recipe, nil
	}

	// if current mineral state is same as asked state, then we are finished.
	if now == stop {
		return recipe, nil
	}

	// get current state order and calculate next state.
	// FIXME: minerals is now lib
	current, _ := minerals.StateTable[now]
	current++

	// apply state rotation
	// FIXME: minerals is now lib
	if current >= len(minerals.StateTable) {
		current = 0
	}

	// convert next state order into real string state.
	nextState, err := minerals.GetByOrder(current)
	if err != nil {
		return recipe, err
	}

	switch nextState {
	case minerals.Fracture:
		recipe = append(recipe, ApplyGrinding)
	case minerals.Liquid:
		recipe = append(recipe, ApplySmelting)
	case minerals.Solid:
		recipe = append(recipe, ApplyFreezing)
	}

	return chainActions(recipe, nextState, stop)
}
