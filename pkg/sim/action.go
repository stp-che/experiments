package sim

type ActionType byte

const (
	AMove ActionType = iota + 1
)

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

// func (a *Action) Changes() []change {
// 	if a.changes == nil {
// 		a.changes = []change{}
// 	}
// 	return a.changes
// }

func (a *Action) IsPossible() bool {
	switch a.Type {
	case AMove:
		return a.targetPos >= 0
	default:
		return false
	}
}

func (a *Action) Apply(ctx map[int]int) {
	if !a.IsPossible() {
		return
	}
	switch a.Type {
	case AMove:
		(&clearReg{a.world.Regions[a.bot.Pos]}).Apply()
		(&putBot{
			Reg: a.world.Regions[a.targetPos],
			Bot: a.bot,
			Pos: a.targetPos,
		}).Apply()
	}
}
