package sim

import (
	"experiments/pkg/sim/behaviour"
	"experiments/pkg/sim/core"
)

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
		b._processResult = b.Brain.Process(b.LookAround(), b.SelfCheck())
	}
	return b._processResult
}

func (b *Bot) IsAlive() bool {
	return b.Energy > 0
}

func (b *Bot) LookAround() []behaviour.OuterInput {
	return []behaviour.OuterInput{
		b.look(core.Up),
		b.look(core.Right),
		b.look(core.Down),
		b.look(core.Left),
	}
}

func (b *Bot) toSignal(c RegionContent) uint8 {
	return uint8(c)
}

func (b *Bot) look(d core.Direction) behaviour.OuterInput {
	r := b.Brain.VisionRange()
	res := behaviour.OuterInput{
		Direction: d,
		Signal:    make([]uint8, r*(r+1)),
	}
	for i := 0; i < len(res.Signal); i++ {
		pos := b.posNear(d, i)
		reg := b.world.Region(pos)
		if reg == nil {
			continue
		}

		res.Signal[i] = b.toSignal(reg.Content)
	}
	return res
}

func (b *Bot) posNear(d core.Direction, i int) int {
	r := b.Brain.VisionRange()
	switch d {
	case core.Up:
		return b.world.ShiftXY(b.Pos, i%(r+1)-r, i/(r+1)-r)
	case core.Right:
		return b.world.ShiftXY(b.Pos, r-i/(r+1), i%(r+1)-r)
	case core.Down:
		return b.world.ShiftXY(b.Pos, r-i%(r+1), r-i/(r+1))
	case core.Left:
		return b.world.ShiftXY(b.Pos, i/(r+1)-r, r-i%(r+1))
	default:
		return -1
	}
}

func (b *Bot) SelfCheck() behaviour.InnerInput {
	return behaviour.InnerInput{b.Energy / 256}
}
