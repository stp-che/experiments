package behaviour

import (
	"math/rand"
)

type IBrain interface {
	Process(OuterInput, InnerInput) *ProcessingResult
}

type Brain struct {
	OuterReceptor       OuterReceptor
	OuterAnalyzersCount int
	OuterAnalyzerNet    OuterAnalyzerNet
	ManipulationSystem  ManipulationSystem
}

func (b *Brain) Process(o OuterInput, i InnerInput) *ProcessingResult {
	collectedSignal := b.OuterReceptor.CollectSignal(b.randomOuterSignal())
	activation := b.OuterAnalyzerNet.Activation(collectedSignal)
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
	b.OuterReceptor = randomOuterReceptor(b.OuterAnalyzersCount)
	b.OuterAnalyzerNet = randomOuterAnalyzerNet(b.OuterAnalyzersCount, len(b.ManipulationSystem))
	return b
}

func (b *Brain) randomOuterSignal() []uint8 {
	size := len(b.OuterReceptor)
	res := make([]uint8, size)
	for i := 0; i < size; i++ {
		res[i] = uint8(rand.Intn(4))
	}
	return res
}
