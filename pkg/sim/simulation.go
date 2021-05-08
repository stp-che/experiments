package sim

type Simulation struct {
	World    *World
	Bot      *Bot
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
		s.Bot.Move()
	}
}

func NewSimulation(cfg Config) *Simulation {
	world := newWorld(cfg.WorldWidth, cfg.WorldHeight)
	bot := &Bot{}
	bot.Attach(world, Pos{0, 0})
	return &Simulation{
		World: world,
		Bot:   bot,
	}
}

type Config struct {
	// Number of regions horizonally
	WorldWidth int
	// Number of regions vertically
	WorldHeight int
	// Number of different genome per simulation
	GenomesNumber int
	// Size of group based on one genome
	GoupSize int
	// Number of mutants per group
	MutantsPerGroup int
}
