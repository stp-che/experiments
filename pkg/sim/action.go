package sim

type ActionType byte

const (
	AMove ActionType = iota + 1
)

var actionTypesByPriority = []ActionType{
	AMove,
}

type Action struct {
	Type      ActionType
	Direction Direction
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

// Checks whether the action is possible considering the current state of the world
// It does not take into account actions that are going to be performed by other bots in this step
func (a *Action) IsPossible() bool {
	switch a.Type {
	case AMove:
		return a.targetPos >= 0 && !a.world.Regions[a.targetPos].Busy()
	default:
		return false
	}
}

// Returns a list of changes caused by the action
// ctx contains number of actions of the same type for the position (i.e. map[targetPos]count)
// So for example for AMove action ctx[targetPos] > 1 means that some other bots are going to move to the same position
// Thus the action can not be performed and has no effect
func (a *Action) Effect(ctx map[int]int) []change {
	if !a.IsPossible() {
		return nil
	}
	switch a.Type {
	case AMove:
		return []change{
			&clearReg{a.world.Regions[a.bot.Pos]},
			&putBot{
				Reg: a.world.Regions[a.targetPos],
				Bot: a.bot,
				Pos: a.targetPos,
			},
		}
	default:
		return nil
	}
}
