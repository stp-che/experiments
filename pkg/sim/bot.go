package sim

type iBotController interface {
	NextAction(*World, int) *Action
}

type randomBotCtrl struct{}

func (c *randomBotCtrl) NextAction(world *World, pos int) *Action {
	return &Action{
		Type:      AMove,
		Direction: randomDirection(),
	}
}

type Bot struct {
	world      *World
	Ctrl       iBotController
	Pos        int
	nextAction *Action
}

func (b *Bot) Init(world *World, ctrl iBotController) *Bot {
	b.world = world
	b.Ctrl = ctrl
	if b.Ctrl == nil {
		b.Ctrl = &randomBotCtrl{}
	}
	return b
}

func (b *Bot) NextAction() *Action {
	if b.nextAction == nil {
		b.nextAction = b.Ctrl.NextAction(b.world, b.Pos).Bind(b.world, b)
	}
	return b.nextAction
}

func (b *Bot) DoNextAction(ctx map[int]int) {
	b.NextAction().Apply(ctx)
	b.nextAction = nil
}

// func (b *Bot) Move() bool {
// 	newPos := b.world.NextPos(b.Pos, randomDirection())
// 	newReg := b.world.Regions[newPos]
// 	if newReg == nil {
// 		return false
// 	}
// 	b.world.Regions[b.Pos].Clear()
// 	b.Pos = newPos
// 	newReg.Occupy(b)
// 	return true
// }
