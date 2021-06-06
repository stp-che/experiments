package sim

import "experiments/pkg/sim/behaviour"

const DEFAULT_ENERGY = 300

type Bot struct {
	world          *World
	Brain          behaviour.IBrain
	Energy         int
	Age            int
	Pos            int
	nextAction     *Action
	_processResult *behaviour.ProcessingResult
}

func (b *Bot) Init(world *World) *Bot {
	b.world = world
	b.Energy = DEFAULT_ENERGY
	return b
}

func (b *Bot) NextAction() *Action {
	d := b.processResult().Decision
	if d != nil && b.nextAction == nil {
		b.nextAction = (&Action{
			Type:      d.ActionType,
			Direction: d.Direction,
		}).Bind(b.world, b)
	}
	return b.nextAction
}

func (b *Bot) StepDone() {
	if !b.IsAlive() {
		return
	}

	energyLost := b.processResult().EnergyCost
	if b.nextAction != nil {
		energyLost *= b.nextAction.EnergyCostMultiplier()
	}
	b.Energy -= energyLost
	b.nextAction = nil
	b._processResult = nil
	b.Age++
}

func (b *Bot) processResult() *behaviour.ProcessingResult {
	if b._processResult == nil {
		b._processResult = b.Brain.Process(behaviour.OuterInput{}, behaviour.InnerInput{})
	}
	return b._processResult
}

func (b *Bot) IsAlive() bool {
	return b.Energy > 0
}
