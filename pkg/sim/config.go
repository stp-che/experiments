package sim

type Config struct {
	// Number of regions horizonally
	WorldWidth int
	// Number of regions vertically
	WorldHeight int
	// Number of different brains per simulation
	BrainsNumber int
	// Size of group based on one genome
	GroupSize int
	// Number of mutants per group
	MutantsPerGroup int
	// Amount of ranges where food appears
	FoodAmount int
}
