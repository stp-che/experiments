package behaviour

import (
	"math/rand"
)

type IBrain interface {
	Process([]OuterInput, InnerInput) *ProcessingResult
	VisionRange() int
	Mutate(int) IBrain
}

type Brain struct {
	OuterReceptor       OuterReceptor
	HealthAnalyzerNet   HealthAnalyzerNet
	OuterAnalyzersCount int
	OuterAnalyzerNet    OuterAnalyzerNet
	ManipulationSystem  ManipulationSystem
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
	return int(b.OuterReceptor.visionRange)
}

func (b *Brain) energyCost(activations ManipulationSystemActivation) int {
	baseCost := b.OuterReceptor.Size() +
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

func (b *Brain) Mutate(n int) IBrain {
	newBrain := b
	for i := 0; i < n; i++ {
		newBrain = newBrain.randomMutation().apply(newBrain)
	}
	return newBrain
}

func (b *Brain) copy() *Brain {
	return &Brain{
		OuterReceptor:       b.OuterReceptor,
		HealthAnalyzerNet:   b.HealthAnalyzerNet,
		OuterAnalyzersCount: b.OuterAnalyzersCount,
		OuterAnalyzerNet:    b.OuterAnalyzerNet,
		ManipulationSystem:  b.ManipulationSystem,
	}
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

func (b *Brain) randomMutation() iMutation {
	switch rand.Intn(6) {
	case 0:
		return mIncreaseVisionRange{}
	case 1:
		return mDecreaseVisionRange{}
	case 2:
		return randomChangeOuterReceptor(b)
	case 3:
		return randomAddHealthAnalyzerLink(b)
	case 4:
		return randomChangeHealthAnalyzerCorrection(b)
	case 5:
		return randomChangeHealthAnalyzerMinMax(b)
	default:
		return randomChangeOuterReceptor(b)
	}
}
