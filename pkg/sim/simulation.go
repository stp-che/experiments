package sim

import (
	"errors"
	"sync"
	"time"
)

type Simulation struct {
	running        bool
	finished       bool
	cfg            Config
	Experiment     *Experiment
	mutex          *sync.Mutex
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
		for !s.Experiment.finished {
			<-s.ticker.C
			s.Experiment.Step()
			updates <- nil
		}
		s.running = false
	}()
	return updates, nil
}

func (s *Simulation) Finished() bool {
	return s.finished
}

func NewSimulation(cfg Config) (*Simulation, error) {
	s := &Simulation{cfg: cfg, mutex: &sync.Mutex{}}
	err := s.init()
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Simulation) init() error {
	s.tickerInterval = time.Millisecond * 128
	ex, err := s.createExperiment()
	s.Experiment = ex
	return err
}

func (s *Simulation) createExperiment() (*Experiment, error) {
	ex := &Experiment{cfg: s.cfg, mutex: s.mutex}
	err := ex.init()
	if err != nil {
		return nil, err
	}
	return ex, nil
}

func (s *Simulation) Sync(f func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	f()
}
