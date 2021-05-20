package sim

type iGenome interface {
	NextAction(*World, int) *Action
}

type Genome struct{}

func (c *Genome) NextAction(world *World, pos int) *Action {
	return randomAction()
}
