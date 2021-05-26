package sim

import "errors"

type BotsGroup struct {
	Genome *Genome
	Bots   []*Bot
}

type Simulation struct {
	cfg          Config
	World        *World
	Bots         []*Bot
	Genomes      []*Genome
	Steps        int
	finished     bool
	foodNextStep int
	groups       []BotsGroup
}

// func (s *Simulation) Run(updates chan interface{}) {
// 	go func() {
// 		for !s.finished {
// 			s.Step()
// 			updates <- nil
// 		}
// 	}()
// }

var (
	NoPlaceForBots = errors.New("No place for bots in the world")
)

func (s *Simulation) Finished() bool {
	return s.finished
}

func NewSimulation(cfg Config) (*Simulation, error) {
	s := &Simulation{}
	err := s.Init(cfg)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Simulation) Init(cfg Config) error {
	s.cfg = cfg
	s.createWorld()
	s.createGenomes()
	err := s.createBots()
	return err
}

func (s *Simulation) createWorld() {
	s.World = newWorld(s.cfg.WorldWidth, s.cfg.WorldHeight)
}

func (s *Simulation) createGenomes() {
	s.Genomes = make([]*Genome, s.cfg.GenomesNumber)
	s.groups = make([]BotsGroup, s.cfg.GenomesNumber)
	for i := 0; i < s.cfg.GenomesNumber; i++ {
		s.Genomes[i] = &Genome{}
		s.groups[i] = BotsGroup{Genome: s.Genomes[i]}
	}
}

func (s *Simulation) createBots() error {
	botsNumber := s.cfg.GenomesNumber * s.cfg.GroupSize
	ps := s.World.RandomEmptyPositions(botsNumber)
	if botsNumber > len(ps) {
		return NoPlaceForBots
	}
	s.Bots = make([]*Bot, botsNumber)
	n := 0
	for i, g := range s.Genomes {
		for j := 0; j < s.cfg.GroupSize; j++ {
			s.Bots[n] = (&Bot{Genome: g}).Init(s.World)
			(&putBot{
				Bot: s.Bots[n],
				Reg: s.World.Regions[ps[n]],
				Pos: ps[n],
			}).Apply()
			n++
		}
		s.groups[i].Bots = s.Bots[i*s.cfg.GroupSize : (i+1)*s.cfg.GroupSize-1]
	}
	return nil
}

func (s *Simulation) Step() {
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

func (s *Simulation) nextActionsAndContexts() (map[ActionType][]*Action, map[ActionType]map[int]int) {
	actionsByType := map[ActionType][]*Action{}
	contexts := map[ActionType]map[int]int{}
	for _, bot := range s.Bots {
		if !bot.IsAlive() {
			continue
		}

		a := bot.NextAction()
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
	for _, pos := range s.World.RandomEmptyPositions(100) {
		s.World.Regions[pos].Content = RCFood
	}
	s.foodNextStep += 30
}

func (s *Simulation) BotsGroups() []BotsGroup {
	return s.groups
}

type Config struct {
	// Number of regions horizonally
	WorldWidth int
	// Number of regions vertically
	WorldHeight int
	// Number of different genome per simulation
	GenomesNumber int
	// Size of group based on one genome
	GroupSize int
	// Number of mutants per group
	MutantsPerGroup int
}
