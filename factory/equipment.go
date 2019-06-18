package factory

import (
	"context"
	"time"
)

// Equipment is some magic mechanism in factory. For example equipment
// can be freezer, grinder, smelter.
type Equipment interface {

	// Empty returns true if equipment is empty. Equipment can hold
	// only one item. Should return true if any mineral is inserted.
	Empty() bool

	// SetPower sets equipment power in watts. 1 watt is used to
	// crack 1 hardness per second or to control 1 temperature per
	// second. Power change while performing action should not
	// affect current action speed.
	SetPower(watts int) error

	// Insert mineral into equipment. Mineral should be inserted into
	// equipment before applying action to this mineral. Equipment
	// must be empty before inserting something.
	Insert(item Mineral) error

	// Takeout should return mineral from equipment. After successful
	// takeout, equipment should become empty. Takeout should fail if
	// action is in progress.
	Takeout() (Mineral, error)

	// Perform action on inserted mineral. Action should change
	// inserted mineral state. After successful perform, mineral can be
	// taked out.
	Perform() error

	// PerformWithCtx is same as Perform, but context should be accepted
	// in order to stop action in middle.
	PerformWithCtx(ctx context.Context) error
}

func CalculateProcessTime(challange, power int) time.Duration {
	// we cannot apply any calculations if power is zero, because there
	// will be divizion by zero.
	if power == 0 {
		return 0
	}
	full := challange / power
	partial := challange % power
	processTime := time.Duration(full) * time.Second
	processTime += time.Duration(partial) * time.Millisecond

	return processTime
}
