package sim

type iGenome interface {
	NextAction(*World, int) *Action
	EnergyCost() int
}

type Genome struct{}

func (g *Genome) NextAction(world *World, pos int) *Action {
	return randomAction()
}

func (g *Genome) EnergyCost() int {
	return 10
}
