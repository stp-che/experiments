package sim

type Simulation struct {
	World    *World
	Bots     []*Bot
	finished bool
}

// func (s *Simulation) Run(updates chan interface{}) {
// 	go func() {
// 		for !s.finished {
// 			s.Step()
// 			updates <- nil
// 		}
// 	}()
// }

func (s *Simulation) Finished() bool {
	return s.finished
}

func NewSimulation(cfg Config) *Simulation {
	world := newWorld(cfg.WorldWidth, cfg.WorldHeight)
	botsNumber := 10
	ps := world.RandomEmptyPositions(botsNumber)
	if botsNumber > len(ps) {
		botsNumber = len(ps)
	}
	bots := createBots(botsNumber, world)
	for i := 0; i < botsNumber; i++ {
		(&putBot{
			Bot: bots[i],
			Reg: world.Regions[ps[i]],
			Pos: ps[i],
		}).Apply()
	}
	return &Simulation{
		World: world,
		Bots:  bots,
	}
}

func createBots(n int, w *World) []*Bot {
	bots := make([]*Bot, n)
	for i := 0; i < n; i++ {
		bots[i] = (&Bot{}).Init(w, nil)
	}
	return bots
}

func (s *Simulation) Step() {
	if !s.finished {
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

		for _, bot := range s.Bots {
			bot.StepDone()
		}
	}
}

func (s *Simulation) nextActionsAndContexts() (map[ActionType][]*Action, map[ActionType]map[int]int) {
	actionsByType := map[ActionType][]*Action{}
	contexts := map[ActionType]map[int]int{}
	for _, bot := range s.Bots {
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
