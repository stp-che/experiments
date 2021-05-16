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

func (b *Bot) StepDone() {
	b.nextAction = nil
}
