package behaviour

type IBrain interface {
	Process([]OuterInput, InnerInput) *ProcessingResult
	VisionRange() int
}

type Brain struct {
	Structure BrainStructure
	Content   []uint8
}

func (b *Brain) Process(o []OuterInput, i InnerInput) *ProcessingResult {
	activation := ManipulationSystemActivation{}
	correction := b.healthAnalyzerNet().Correction(i[0])
	for _, inp := range o {
		collectedSignal := b.outerReceptor().CollectSignal(inp.Signal)
		activation[inp.Direction] = b.outerAnalyzerNet().Activation(collectedSignal, correction)
	}
	return &ProcessingResult{
		Decision:   b.manipulationSystem().ComputeIntention(activation),
		EnergyCost: b.energyCost(activation),
	}
}

func (b *Brain) VisionRange() int {
	return int(b.Structure.VisionRange)
}

func (b *Brain) outerReceptor() OuterReceptor {
	return OuterReceptor(b.Content[b.Structure.outerReceptorStart():b.Structure.outerReceptorEnd()])
}

func (b *Brain) healthAnalyzerNet() HealthAnalyzerNet {
	return HealthAnalyzerNet(b.Content[b.Structure.healthAnalyzerNetStart():b.Structure.healthAnalyzerNetEnd()])
}

func (b *Brain) outerAnalyzerNet() OuterAnalyzerNet {
	return OuterAnalyzerNet(b.Content[b.Structure.outerAnalyzerNetStart():b.Structure.outerAnalyzerNetEnd()])
}

func (b *Brain) manipulationSystem() ManipulationSystem {
	return ManipulationSystem(b.Content[b.Structure.manipulationSystemStart():b.Structure.manipulationSystemEnd()])
}

func (b *Brain) energyCost(activations ManipulationSystemActivation) int {
	baseCost := len(b.Content) / 4
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

func (b *Brain) Mutate(n int) *Brain {
	return nil
}

func (b *Brain) NormalizeContent() {
	b.outerReceptor().normalize(b.Structure)
	b.healthAnalyzerNet().normalize(b.Structure)
	b.outerAnalyzerNet().normalize(b.Structure)
	b.manipulationSystem().normalize(b.Structure)
}

func RandomBrain() *Brain {
	return RandomBrainWithStructure(randomBrainStructure())
}

func RandomBrainWithStructure(s BrainStructure) *Brain {
	b := &Brain{
		Structure: s,
		Content:   make([]uint8, s.contentSize()),
	}
	b.outerReceptor().randomize(b.Structure)
	b.healthAnalyzerNet().randomize(b.Structure)
	b.outerAnalyzerNet().randomize(b.Structure)
	b.manipulationSystem().randomize(b.Structure)
	return b
}
