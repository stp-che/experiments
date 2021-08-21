package behaviour

import (
	"math/rand"
)

const (
	maxOuterAnalyzerCount = 255
	maxManipulatorsCount  = 255
)

type reductiveMutation struct{}

func (m reductiveMutation) apply(b *Brain) *Brain {
	newBrain := b.copy()

	idx := findOrphanManipulator(newBrain)
	if idx >= 0 {
		removeManipulator(newBrain, idx)
		return newBrain
	}

	idx = findUselessOuterAnalyzer(newBrain)
	if idx >= 0 {
		removeOuterAnalyzer(newBrain, uint8(idx))
		return newBrain
	}

	idx = findUselessOuterLink(newBrain)
	if idx >= 0 {
		removeOuterAnalyzerLink(newBrain, idx)
		return newBrain
	}

	idx = findUselessHealthLink(newBrain)
	if idx >= 0 {
		removeHealthAnalyzerLink(newBrain, idx)
		return newBrain
	}

	return newBrain
}

func findOrphanManipulator(b *Brain) int {
	refs := make(map[uint8]bool)
	for _, link := range b.OuterAnalyzerNet {
		refs[link.Manipulator] = true
	}
	size := uint8(len(b.ManipulationSystem))
	var i uint8
	for ; i < size; i++ {
		if !refs[i] {
			return int(i)
		}
	}
	return -1
}

func removeManipulator(b *Brain, idx int) {
	newSys := make(ManipulationSystem, len(b.ManipulationSystem)-1)
	for i, manipulator := range b.ManipulationSystem {
		if i < idx {
			newSys[i] = manipulator
		} else if i > idx {
			newSys[i-1] = manipulator
		}
	}
	b.ManipulationSystem = newSys
	newNet := b.OuterAnalyzerNet.copy()
	for i, link := range newNet {
		if link.Manipulator > uint8(idx) {
			newLink := link.copy()
			newLink.Manipulator--
			newNet[i] = &newLink
		}
	}
	b.OuterAnalyzerNet = newNet
}

func findUselessOuterLink(b *Brain) int {
	refs := make(map[uint8]bool)
	for _, link := range b.HealthAnalyzerNet {
		refs[link.OuterAnalyzerLink] = true
	}
	for i, link := range b.OuterAnalyzerNet {
		if link.Power == 0 && !refs[uint8(i)] {
			return int(i)
		}
	}
	return -1
}

func findUselessOuterAnalyzer(b *Brain) int {
	refs := make(map[uint8]bool)
	for _, link := range b.OuterAnalyzerNet {
		refs[link.Analyzer] = true
	}
	size := uint8(b.OuterAnalyzersCount)
	var i uint8
	for ; i < size; i++ {
		if !refs[i] {
			return int(i)
		}
	}
	return -1
}

func removeOuterAnalyzer(b *Brain, analyzer uint8) {
	b.OuterAnalyzersCount--

	newRec := NewOuterReceptor(b.OuterReceptor.visionRange)
	for i, row := range b.OuterReceptor.cells {
		for j, a := range row {
			if a > analyzer || a > uint8(b.OuterAnalyzersCount) {
				newRec.cells[i][j] = a - 1
			} else {
				newRec.cells[i][j] = a
			}
		}
	}

	b.OuterReceptor = newRec
}

func findUselessHealthLink(b *Brain) int {
	for i, link := range b.HealthAnalyzerNet {
		if b.OuterAnalyzerNet[int(link.OuterAnalyzerLink)].Power == 0 {
			return i
		}
	}
	return -1
}

type mAddOuterAnalyzer struct{}

func (m mAddOuterAnalyzer) apply(b *Brain) *Brain {
	newBrain := b.copy()
	if newBrain.OuterAnalyzersCount < maxOuterAnalyzerCount {
		newBrain.OuterAnalyzersCount++
	}
	return newBrain
}

type mAddManipulator struct {
	action    ActionType
	dirValues [8]int8
}

func (m mAddManipulator) apply(b *Brain) *Brain {
	newBrain := b.copy()
	if len(b.ManipulationSystem) < maxManipulatorsCount {
		newManipulator := Manipulator{m.action, m.dirValues}
		newBrain.ManipulationSystem = append(b.ManipulationSystem, &newManipulator)
	}
	return newBrain
}

func randomAddManipulatorMutation() mAddManipulator {
	dirValues := [8]int8{}
	dirValues[rand.Intn(8)] = int8(rand.Intn(manipulatorValueMaxDelta*2+1) - manipulatorValueMaxDelta)
	return mAddManipulator{
		action:    randomActionType(),
		dirValues: dirValues,
	}
}
