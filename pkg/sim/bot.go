package sim

const DEFAULT_ENERGY = 300

type Bot struct {
	world              *World
	Genome             iGenome
	Energy             int
	Age                int
	Pos                int
	nextAction         *Action
	nextActionComputed bool
}

func (b *Bot) Init(world *World) *Bot {
	b.world = world
	b.Energy = DEFAULT_ENERGY
	return b
}

func (b *Bot) NextAction() *Action {
	if !b.nextActionComputed {
		b.nextAction = b.Genome.NextAction(b.world, b.Pos)
		if b.nextAction != nil {
			b.nextAction.Bind(b.world, b)
		}
		b.nextActionComputed = true
	}
	return b.nextAction
}

func (b *Bot) StepDone() {
	if !b.IsAlive() {
		return
	}

	energyLost := b.Genome.EnergyCost()
	if b.nextAction != nil {
		energyLost *= b.nextAction.EnergyCostMultiplier()
	}
	b.Energy -= energyLost
	b.nextAction = nil
	b.nextActionComputed = false
	b.Age++
}

func (b *Bot) IsAlive() bool {
	return b.Energy > 0
}
