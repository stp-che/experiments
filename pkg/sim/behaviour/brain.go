package behaviour

import (
	"experiments/pkg/sim/core"
	"math/rand"
)

type OuterInput struct {
	Direction core.Direction
	Signal    []uint8
}

type InnerInput []int

type ProcessingResult struct {
	Decision   *Intention
	EnergyCost int
}

type IBrain interface {
	Process(OuterInput, InnerInput) *ProcessingResult
}

type Brain struct {
	OuterAnalyzersCount int
	OuterAnalyzerNet    OuterAnalyzerNet
	ManipulationSystem  ManipulationSystem
}

func (b *Brain) Process(o OuterInput, i InnerInput) *ProcessingResult {
	activation := b.OuterAnalyzerNet.Activation(b.randomOuterSignal())
	return &ProcessingResult{
		Decision:   b.ManipulationSystem.ComputeIntention(activation),
		EnergyCost: 10,
	}
}

func RandomBrain() *Brain {
	b := &Brain{
		OuterAnalyzersCount: rand.Intn(10) + 1,
		ManipulationSystem:  randomManipulationSystem(),
	}
	b.OuterAnalyzerNet = randomOuterAnalyzerNet(b.OuterAnalyzersCount, len(b.ManipulationSystem))
	return b
}

func (b *Brain) randomOuterSignal() map[uint8][]uint8 {
	res := make(map[uint8][]uint8)
	for i := 0; i < 30; i++ {
		analyser := rand.Intn(b.OuterAnalyzersCount)
		if _, ok := res[uint8(analyser)]; !ok {
			res[uint8(analyser)] = make([]uint8, 4)
		}
		res[uint8(analyser)][rand.Intn(4)] += 1
	}
	return res
}
