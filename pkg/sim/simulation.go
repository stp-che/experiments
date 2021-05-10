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

func (s *Simulation) Step() {
	if !s.finished {
		for _, b := range s.Bots {
			b.DoNextAction(map[int]int{})
		}
	}
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
