package behaviour

import "math/rand"

type brainPartsBounds struct {
	outerReceptorStart      int
	healthAnalyzerNetStart  int
	outerAnalyzerNetStart   int
	manipulationSystemStart int
	manipulationSystemEnd   int
}

type BrainStructure struct {
	VisionRange            uint8
	OuterAnalyzersCount    uint8
	HealthAnalyzerNetSize  uint8
	OuterAnalyzerNetSize   uint8
	ManipulationSystemSize uint8
	cachedBounds           brainPartsBounds
}

func (s BrainStructure) outerReceptorStart() int {
	return s.bounds().outerReceptorStart
}

func (s BrainStructure) outerReceptorEnd() int {
	return s.bounds().healthAnalyzerNetStart
}

func (s BrainStructure) healthAnalyzerNetStart() int {
	return s.bounds().healthAnalyzerNetStart
}

func (s BrainStructure) healthAnalyzerNetEnd() int {
	return s.bounds().outerAnalyzerNetStart
}

func (s BrainStructure) outerAnalyzerNetStart() int {
	return s.bounds().outerAnalyzerNetStart
}

func (s BrainStructure) outerAnalyzerNetEnd() int {
	return s.bounds().manipulationSystemStart
}

func (s BrainStructure) manipulationSystemStart() int {
	return s.bounds().manipulationSystemStart
}

func (s BrainStructure) manipulationSystemEnd() int {
	return s.bounds().manipulationSystemEnd
}

func (s BrainStructure) contentSize() int {
	return s.bounds().manipulationSystemEnd
}

func (s BrainStructure) bounds() brainPartsBounds {
	b := s.cachedBounds
	if b.manipulationSystemEnd == 0 {
		vr := int(s.VisionRange)
		b.healthAnalyzerNetStart = vr * (vr + 1)
		b.outerAnalyzerNetStart = b.healthAnalyzerNetStart + int(s.HealthAnalyzerNetSize)*healthAnalyzerLinkSize
		b.manipulationSystemStart = b.outerAnalyzerNetStart + int(s.OuterAnalyzerNetSize)*outerAnalyzerLinkSize
		b.manipulationSystemEnd = b.manipulationSystemStart + int(s.ManipulationSystemSize)*manipulatorSize
	}
	return b
}

func randomBrainStructure() BrainStructure {
	s := BrainStructure{
		VisionRange:         uint8(rand.Intn(15)),
		OuterAnalyzersCount: uint8(rand.Intn(10) + 1),
	}
	s.ManipulationSystemSize = uint8(rand.Intn(5) + 1)
	s.OuterAnalyzerNetSize = uint8(rand.Intn(int(s.OuterAnalyzersCount))*rand.Intn(int(s.ManipulationSystemSize)) + 1)
	s.HealthAnalyzerNetSize = uint8(rand.Intn(int(s.OuterAnalyzerNetSize+1)/2) + 1)
	return s
}
