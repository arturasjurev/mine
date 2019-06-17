package factory

// Possible Mineral states. State transaction:
// Fracture -> Dust -> Liquid -> Solid.
const (
	// Fracture is first mineral state. Unpure mineral.
	Fracture = "fracture"

	// When we put fracture mineral into grinder it becomes dust.
	Dust = "dust"

	// When we put dust in smelter it becomes liquid.
	Liquid = "liquid"

	// When we put liquid into freezer it becomes solid.
	Solid = "solid"
)

type MineralState string

type Mineral struct {
	Name         string       `json:"name"`
	MeltingPoint int          `json:"melting_point"`
	Hardness     int          `json:"hardness"`
	Fractures    int          `json:"fractures"`
	State        MineralState `json:"state"`
}
