package behaviour

import (
	"math"
	"math/rand"
)

type IBrain interface {
	Process([]OuterInput, InnerInput) *ProcessingResult
	VisionRange() int
}

type Brain struct {
	OuterReceptor       OuterReceptor
	OuterAnalyzersCount int
	OuterAnalyzerNet    OuterAnalyzerNet
	ManipulationSystem  ManipulationSystem
	visionRange         int
}

func (b *Brain) Process(o []OuterInput, i InnerInput) *ProcessingResult {
	activation := ManipulationSystemActivation{}
	for _, inp := range o {
		collectedSignal := b.OuterReceptor.CollectSignal(inp.Signal)
		activation[inp.Direction] = b.OuterAnalyzerNet.Activation(collectedSignal)
	}
	return &ProcessingResult{
		Decision:   b.ManipulationSystem.ComputeIntention(activation),
		EnergyCost: 10,
	}
}

func (b *Brain) VisionRange() int {
	if b.OuterReceptor == nil {
		return 0
	}
	if b.visionRange > 0 {
		return b.visionRange
	}
	b.visionRange = (int(math.Round(math.Sqrt(float64(4*len(b.OuterReceptor)+1)))) - 1) / 2
	return b.visionRange
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
