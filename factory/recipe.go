package factory

import (
	"fmt"
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

// Define mineral state transformation line. More information can be found on
// facotry/mineral.go file. This table is used to determinate what next actions
// should be performed on mineral in what order.
var stateTable = map[MineralState]int{
	Fracture: 0,
	Liquid:   1,
	Solid:    2,
}

func getByOrder(order int) (MineralState, error) {
	for i, v := range stateTable {
		if v == order {
			return i, nil
		}
	}
	return "", fmt.Errorf("unknown mineral order")
}

// GenerateRecipe will define action order what equipments should be used and
// what actions should be applied to current mineral to reach asked mineral state.
// E.g.: if we pass current = fracture, asked = solid, this function should return
// recipe: []RecipeAction{"apply_grinding", "apply_smelting", "apply_freezing"}
func GenerateRecipe(current, asked MineralState) ([]RecipeAction, error) {

	recipe := []RecipeAction{}

	// check if current and asked states is know for us. If states are
	// unknown for us, we do not know how to reach those states.
	if _, ok := stateTable[current]; !ok {
		return recipe, fmt.Errorf("current state '%s' is unknown", current)
	}
	if _, ok := stateTable[asked]; !ok {
		return recipe, fmt.Errorf("asked state '%s' in unknown", asked)
	}

	return chainActions(recipe, current, asked)
}

// chainActions is core logic how to chain recipe actions. This will return
// recipe in form []RecipeAction{"apply_grinding", "apply_smelting", "apply_freezing"}.
// This function is used inside GenerateRecipe function.
// SPOILER ALERT: RECURSION USED
func chainActions(recipe []RecipeAction, now, stop MineralState) ([]RecipeAction, error) {

	// NOTICE: there is special case. Only on fractured state it is
	// possible to apply grinding again, because sometimes what we
	// want to achieve is to double its fractures. Because of this
	// we can process order like `order{from: fractured, to: fractured}`
	// in case we want to double its fractures.
	if now == stop && stop == Fracture && len(recipe) == 0 {
		recipe = append(recipe, ApplyGrinding)
		return recipe, nil
	}

	// if current mineral state is same as asked state, then we are finished.
	if now == stop {
		return recipe, nil
	}

	// get current state order and calculate next state.
	current, _ := stateTable[now]
	current++

	// apply state rotation
	if current >= len(stateTable) {
		current = 0
	}

	// convert next state order into real string state.
	nextState, err := getByOrder(current)
	if err != nil {
		return recipe, err
	}

	switch nextState {
	case Fracture:
		recipe = append(recipe, ApplyGrinding)
	case Liquid:
		recipe = append(recipe, ApplySmelting)
	case Solid:
		recipe = append(recipe, ApplyFreezing)
	}

	return chainActions(recipe, nextState, stop)
}
