package sim

type Bot struct {
	world      *World
	Genome     iGenome
	Pos        int
	nextAction *Action
}

func (b *Bot) Init(world *World) *Bot {
	b.world = world
	return b
}

func (b *Bot) NextAction() *Action {
	if b.nextAction == nil {
		b.nextAction = b.Genome.NextAction(b.world, b.Pos).Bind(b.world, b)
	}
	return b.nextAction
}

func (b *Bot) StepDone() {
	b.nextAction = nil
}
