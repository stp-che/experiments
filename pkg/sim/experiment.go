package sim

import (
	"experiments/pkg/sim/behaviour"
	"sort"
	"sync"
)

type BotsGroup struct {
	Brain *behaviour.Brain
	Bots  []*Bot
}

type Experiment struct {
	finished     bool
	cfg          Config
	mutex        *sync.Mutex
	Number       int
	World        *World
	Bots         []*Bot
	Brains       []*behaviour.Brain
	Steps        int
	foodNextStep int
	groups       []BotsGroup
}

func (e *Experiment) init() error {
	e.createWorld()
	e.createGenomes()
	err := e.createBots()
	return err
}

func (e *Experiment) createWorld() {
	e.World = newWorld(e.cfg.WorldWidth, e.cfg.WorldHeight)
}

func (e *Experiment) createGenomes() {
	e.Brains = make([]*behaviour.Brain, e.cfg.BrainsNumber)
	e.groups = make([]BotsGroup, e.cfg.BrainsNumber)
	for i := 0; i < e.cfg.BrainsNumber; i++ {
		e.Brains[i] = behaviour.RandomBrain()
		e.groups[i] = BotsGroup{Brain: e.Brains[i]}
	}
}

func (e *Experiment) createBots() error {
	botsNumber := e.cfg.BrainsNumber * e.cfg.GroupSize
	ps := e.World.RandomEmptyPositions(botsNumber)
	if botsNumber > len(ps) {
		return NoPlaceForBots
	}
	e.Bots = make([]*Bot, botsNumber)
	n := 0
	for i, b := range e.Brains {
		for j := 0; j < e.cfg.GroupSize; j++ {
			e.Bots[n] = (&Bot{Brain: b}).Init(e.World)
			(&putBot{
				Bot: e.Bots[n],
				Reg: e.World.Regions[ps[n]],
				Pos: ps[n],
			}).Apply()
			n++
		}
		e.groups[i].Bots = e.Bots[i*e.cfg.GroupSize : (i+1)*e.cfg.GroupSize]
	}
	return nil
}

func (e *Experiment) BotsGroups() []BotsGroup {
	return e.groups
}

func (e *Experiment) Step() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if !e.finished {
		if e.Steps == e.foodNextStep {
			e.generateFood()
		}
		e.performStepActions()
		aliveBotsCount := e.updateBotsStates()
		e.Steps++
		e.finished = e.finished || aliveBotsCount == 0
	}
}

func (e *Experiment) performStepActions() {
	actionsByType, contexts := e.nextActionsAndContexts()
	for _, aType := range actionTypesByPriority {
		actions, ok := actionsByType[aType]
		if !ok {
			continue
		}
		changes := make([]change, 0, 4)
		for _, a := range actions {
			actionChanges := a.Effect(contexts[aType])
			if actionChanges == nil {
				continue
			}
			changes = append(changes, actionChanges...)
		}
		for _, c := range changes {
			c.Apply()
		}
	}
}

func (e *Experiment) nextActionsAndContexts() (map[behaviour.ActionType][]*Action, map[behaviour.ActionType]map[int]int) {
	actionsByType := map[behaviour.ActionType][]*Action{}
	contexts := map[behaviour.ActionType]map[int]int{}
	for _, bot := range e.Bots {
		if !bot.IsAlive() {
			continue
		}

		a := bot.NextAction()
		if a == nil {
			continue
		}

		if actions, ok := actionsByType[a.Type]; ok {
			actionsByType[a.Type] = append(actions, a)
		} else {
			actionsByType[a.Type] = make([]*Action, 1, 4)
			actionsByType[a.Type][0] = a
		}
		if ctx, ok := contexts[a.Type]; !ok {
			contexts[a.Type] = map[int]int{a.TargetPos(): 1}
		} else {
			if _, ok = ctx[a.TargetPos()]; ok {
				ctx[a.TargetPos()] += 1
			} else {
				ctx[a.TargetPos()] = 1
			}
		}
	}

	return actionsByType, contexts
}

func (e *Experiment) updateBotsStates() int {
	aliveBotsCount := 0
	for _, bot := range e.Bots {
		if !bot.IsAlive() {
			continue
		}

		aliveBotsCount++
		bot.StepDone()
		if !bot.IsAlive() {
			(&clearReg{e.World.Regions[bot.Pos]}).Apply()
		}
	}
	return aliveBotsCount
}

func (e *Experiment) generateFood() {
	for _, pos := range e.World.RandomEmptyPositions(e.cfg.FoodAmount) {
		e.World.Regions[pos].Content = RCFood
	}
	e.foodNextStep += 30
}

func (e *Experiment) BotsChart() []*Bot {
	bots := make([]*Bot, len(e.Bots))
	copy(bots, e.Bots)
	sort.Slice(bots, func(i, j int) bool {
		b1, b2 := bots[i], bots[j]
		if b1.Movements < 2 && b2.Movements >= 2 {
			return false
		}
		if b1.Movements >= 2 && b2.Movements < 2 {
			return true
		}
		if b1.Age == b2.Age {
			return b1.Energy > b2.Energy
		}
		return b1.Age > b2.Age
	})
	return bots[0:8]
}
