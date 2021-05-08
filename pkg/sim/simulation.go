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
			b.Move()
		}
	}
}

func NewSimulation(cfg Config) *Simulation {
	world := newWorld(cfg.WorldWidth, cfg.WorldHeight)
	bots := createBots(10)
	for _, b := range bots {
		b.Settle(world, world.RandomPos())
	}
	return &Simulation{
		World: world,
		Bots:  bots,
	}
}

func createBots(n int) []*Bot {
	bots := make([]*Bot, n)
	for i := 0; i < n; i++ {
		bots[i] = &Bot{}
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
