package sim

import (
	"errors"
	"experiments/pkg/sim/behaviour"
	"sort"
	"sync"
	"time"
)

type BotsGroup struct {
	Brain *behaviour.Brain
	Bots  []*Bot
}

type Simulation struct {
	running        bool
	finished       bool
	cfg            Config
	Epoch          int
	World          *World
	Bots           []*Bot
	Brains         []*behaviour.Brain
	Steps          int
	foodNextStep   int
	groups         []BotsGroup
	mutex          sync.Mutex
	tickerInterval time.Duration
	ticker         *time.Ticker
}

var (
	AlreadyRunningOrFinished = errors.New("Simulation is already running or finished")
	NoPlaceForBots           = errors.New("No place for bots in the world")
)

func (s *Simulation) Run() (chan interface{}, error) {
	s.mutex.Lock()
	if s.running || s.finished {
		return nil, AlreadyRunningOrFinished
	}
	s.running = true
	s.ticker = time.NewTicker(s.tickerInterval)
	s.mutex.Unlock()

	updates := make(chan interface{})
	go func() {
		for !s.finished {
			<-s.ticker.C
			s.Step()
			updates <- nil
		}
	}()
	return updates, nil
}

func (s *Simulation) Finished() bool {
	return s.finished
}

func NewSimulation(cfg Config) (*Simulation, error) {
	s := &Simulation{cfg: cfg}
	err := s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Simulation) init() error {
	s.tickerInterval = time.Millisecond * 1024
	s.createWorld()
	s.createGenomes()
	err := s.createBots()
	return err
}

func (s *Simulation) createWorld() {
	s.World = newWorld(s.cfg.WorldWidth, s.cfg.WorldHeight)
}

func (s *Simulation) createGenomes() {
	s.Brains = make([]*behaviour.Brain, s.cfg.BrainsNumber)
	s.groups = make([]BotsGroup, s.cfg.BrainsNumber)
	for i := 0; i < s.cfg.BrainsNumber; i++ {
		s.Brains[i] = behaviour.RandomBrain()
		s.groups[i] = BotsGroup{Brain: s.Brains[i]}
	}
}

func (s *Simulation) createBots() error {
	botsNumber := s.cfg.BrainsNumber * s.cfg.GroupSize
	ps := s.World.RandomEmptyPositions(botsNumber)
	if botsNumber > len(ps) {
		return NoPlaceForBots
	}
	s.Bots = make([]*Bot, botsNumber)
	n := 0
	for i, b := range s.Brains {
		for j := 0; j < s.cfg.GroupSize; j++ {
			s.Bots[n] = (&Bot{Brain: b}).Init(s.World)
			(&putBot{
				Bot: s.Bots[n],
				Reg: s.World.Regions[ps[n]],
				Pos: ps[n],
			}).Apply()
			n++
		}
		s.groups[i].Bots = s.Bots[i*s.cfg.GroupSize : (i+1)*s.cfg.GroupSize]
	}
	return nil
}

func (s *Simulation) Step() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if !s.finished {
		if s.Steps == s.foodNextStep {
			s.generateFood()
		}
		s.performStepActions()
		aliveBotsCount := s.updateBotsStates()
		s.Steps++
		s.finished = s.finished || aliveBotsCount == 0
	}
}

func (s *Simulation) Sync(f func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	f()
}

func (s *Simulation) performStepActions() {
	actionsByType, contexts := s.nextActionsAndContexts()
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

func (s *Simulation) nextActionsAndContexts() (map[behaviour.ActionType][]*Action, map[behaviour.ActionType]map[int]int) {
	actionsByType := map[behaviour.ActionType][]*Action{}
	contexts := map[behaviour.ActionType]map[int]int{}
	for _, bot := range s.Bots {
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

func (s *Simulation) updateBotsStates() int {
	aliveBotsCount := 0
	for _, bot := range s.Bots {
		if !bot.IsAlive() {
			continue
		}

		aliveBotsCount++
		bot.StepDone()
		if !bot.IsAlive() {
			(&clearReg{s.World.Regions[bot.Pos]}).Apply()
		}
	}
	return aliveBotsCount
}

func (s *Simulation) generateFood() {
	for _, pos := range s.World.RandomEmptyPositions(s.cfg.FoodAmount) {
		s.World.Regions[pos].Content = RCFood
	}
	s.foodNextStep += 30
}

func (s *Simulation) BotsGroups() []BotsGroup {
	return s.groups
}

func (s *Simulation) BotsChart() []*Bot {
	bots := make([]*Bot, len(s.Bots))
	copy(bots, s.Bots)
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

type Config struct {
	// Number of regions horizonally
	WorldWidth int
	// Number of regions vertically
	WorldHeight int
	// Number of different brains per simulation
	BrainsNumber int
	// Size of group based on one genome
	GroupSize int
	// Number of mutants per group
	MutantsPerGroup int
	// Amount of ranges where food appears
	FoodAmount int
}
