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
	HealthAnalyzerNet   HealthAnalyzerNet
	OuterAnalyzersCount int
	OuterAnalyzerNet    OuterAnalyzerNet
	ManipulationSystem  ManipulationSystem
	visionRange         int
}

func (b *Brain) Process(o []OuterInput, i InnerInput) *ProcessingResult {
	activation := ManipulationSystemActivation{}
	correction := b.HealthAnalyzerNet.Correction(i[0])
	for _, inp := range o {
		collectedSignal := b.OuterReceptor.CollectSignal(inp.Signal)
		activation[inp.Direction] = b.OuterAnalyzerNet.Activation(collectedSignal, correction)
	}
	return &ProcessingResult{
		Decision:   b.ManipulationSystem.ComputeIntention(activation),
		EnergyCost: b.energyCost(activation),
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

func (b *Brain) energyCost(activations ManipulationSystemActivation) int {
	baseCost := len(b.OuterReceptor) +
		len(b.HealthAnalyzerNet) +
		len(b.OuterAnalyzerNet) +
		len(b.ManipulationSystem) +
		b.OuterAnalyzersCount + 1 // number of outer analyzers and one inner analyzer
	activityCost := 0
	for _, activation := range activations {
		for _, pow := range activation {
			if pow < 0 {
				pow = -pow
			}
			activityCost += int(pow)
		}
	}
	return baseCost + (activityCost-1)/32 + 1
}

func RandomBrain() *Brain {
	b := &Brain{
		OuterAnalyzersCount: rand.Intn(10) + 1,
		ManipulationSystem:  randomManipulationSystem(),
	}
	b.OuterReceptor = randomOuterReceptor(b.OuterAnalyzersCount)
	b.OuterAnalyzerNet = randomOuterAnalyzerNet(b.OuterAnalyzersCount, len(b.ManipulationSystem))
	b.HealthAnalyzerNet = randomHealthAnalyzerNet(len(b.OuterAnalyzerNet))
	return b
}
