package factory

import (
	"context"
	"time"

	"github.com/sheirys/mine/minerals"
)

// Equipment is some magic mechanism in factory. For example equipment
// can be freezer, grinder, smelter.
type Equipment interface {

	// Empty returns true if equipment is empty. Equipment can hold
	// only one item. Should return true if any mineral is inserted.
	Empty() bool

	// SetPower sets equipment power in watts. 1 watt is used to
	// crack 1 hardness per second or to control 1 temperature per
	// second. Power change while processing should not affect
	// current action speed.
	SetPower(watts int) error

	// Insert mineral into equipment. Mineral should be inserted into
	// equipment before applying process to this mineral. Equipment
	// must be empty before inserting something.
	Insert(item minerals.Mineral) error

	// Takeout should return mineral from equipment. After successful
	// takeout, equipment should become empty. Takeout should fail if
	// process is in progress.
	Takeout() (minerals.Mineral, error)

	// Process action on inserted mineral. Action should change
	// inserted mineral state. After successful process, mineral can be
	// taked out.
	Process() error

	// ProcessWithCtx is same as Process, but context should be accepted
	// in order to stop action in middle.
	ProcessWithCtx(ctx context.Context) error
}

// CalculateProcessTime will calculate how long equipment should process inserted
// mineral to change its state. Here, depending on equipment, mineral hardness or
// mineral melting point should be passed as challenge and equipment power as power.
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
