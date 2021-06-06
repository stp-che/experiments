package sim

import (
	"experiments/pkg/sim/behaviour"
	"experiments/pkg/sim/core"
)

var actionTypesByPriority = []behaviour.ActionType{
	behaviour.AEat,
	behaviour.AMove,
}

var energyCostMultipliers = map[behaviour.ActionType]int{
	behaviour.AMove: 2,
}

type Action struct {
	Type      behaviour.ActionType
	Direction core.Direction
	world     *World
	bot       *Bot
	targetPos int
}

func (a *Action) Bind(w *World, b *Bot) *Action {
	a.world = w
	a.bot = b
	a.targetPos = w.NextPos(b.Pos, a.Direction)
	return a
}

func (a *Action) TargetPos() int {
	return a.targetPos
}

func (a *Action) TargetReg() *Region {
	return a.world.Regions[a.targetPos]
}

// Checks whether the action is possible considering the current state of the world
// It does not take into account actions that are going to be performed by other bots in this step
func (a *Action) IsPossible() bool {
	if a.targetPos < 0 {
		return false
	}
	switch a.Type {
	case behaviour.AMove:
		return !a.TargetReg().Busy()
	case behaviour.AEat:
		return a.TargetReg().Content == RCFood
	default:
		return false
	}
}

func (a *Action) hasEffect(ctx map[int]int) bool {
	if !a.IsPossible() {
		return false
	}
	switch a.Type {
	case behaviour.AMove:
		return ctx[a.targetPos] == 1
	default:
		return true
	}
}

// Returns a list of changes caused by the action
// ctx contains number of actions of the same type for the position (i.e. map[targetPos]count)
// So for example for AMove action ctx[targetPos] > 1 means that some other bots are going to move to the same position
// Thus the action can not be performed and has no effect
func (a *Action) Effect(ctx map[int]int) []change {
	if !a.hasEffect(ctx) {
		return nil
	}
	switch a.Type {
	case behaviour.AMove:
		return []change{
			&clearReg{a.world.Regions[a.bot.Pos]},
			&putBot{
				Reg: a.TargetReg(),
				Bot: a.bot,
				Pos: a.targetPos,
			},
		}
	case behaviour.AEat:
		return []change{
			&feedBot{
				Bot: a.bot,
				// food is shared equally among all bots eating from the same region
				Energy: 200 / ctx[a.targetPos],
			},
			&clearReg{a.TargetReg()},
		}
	default:
		return nil
	}
}

func (a *Action) EnergyCostMultiplier() int {
	if m, set := energyCostMultipliers[a.Type]; set {
		return m
	}
	return 1
}

// func randomAction() *Action {
// 	return &Action{
// 		Type:      randomActionType(),
// 		Direction: core.RandomDirection(),
// 	}
// }

// func randomActionType() behaviour.ActionType {
// 	return actionTypesByPriority[rand.Intn(len(actionTypesByPriority))]
// }
