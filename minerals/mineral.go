package minerals

import "fmt"

// Possible Mineral states. Possible mineral transactions:
// fractured -> liquid.
// liquid -> solid.
// solid -> fractured.
const (
	Fracture = "fractured"
	Liquid   = "liquid"
	Solid    = "solid"
)

type State string

type Mineral struct {
	Name         string `json:"name"`
	State        State  `json:"state"`
	MeltingPoint int    `json:"melting_point"`
	Hardness     int    `json:"hardness"`
	Fractures    int    `json:"fractures"`
}

// Define mineral state transformation line. More information can be found on
// facotry/mineral.go file. This table is used to determinate what next actions
// should be performed on mineral in what order.
var StateTable = map[State]int{
	Fracture: 0,
	Liquid:   1,
	Solid:    2,
}

func GetByOrder(order int) (State, error) {
	for i, v := range StateTable {
		if v == order {
			return i, nil
		}
	}
	return "", fmt.Errorf("unknown mineral order")
}
