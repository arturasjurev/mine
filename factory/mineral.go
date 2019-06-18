package factory

// Possible Mineral states. Possible mineral transactions:
// fractured -> liquid.
// liquid -> solid.
// solid -> fractured.
const (
	Fracture = "fractured"
	Liquid   = "liquid"
	Solid    = "solid"
)

type MineralState string

type Mineral struct {
	Name         string       `json:"name"`
	State        MineralState `json:"state"`
	MeltingPoint int          `json:"melting_point"`
	Hardness     int          `json:"hardness"`
	Fractures    int          `json:"fractures"`
}
