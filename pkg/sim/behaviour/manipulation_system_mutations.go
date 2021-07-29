package behaviour

import (
	"experiments/pkg/sim/core"
	"math/rand"
)

const manipulatorValueMaxDelta = 10

type mChangeManipulatorValue struct {
	manipulator int
	direction   core.Direction
	delta       int
}

func (m mChangeManipulatorValue) apply(brain *Brain) *Brain {
	newBrain := brain.copy()
	if m.manipulator < len(brain.ManipulationSystem) && m.direction.IsValid() {
		newBrain.ManipulationSystem = newBrain.ManipulationSystem.copy()
		newManipulator := newBrain.ManipulationSystem[m.manipulator].copy()
		newValue := cutToInt8(int(newManipulator.DirValues[m.direction-1]) + m.delta)
		newManipulator.DirValues[m.direction-1] = newValue
		newBrain.ManipulationSystem[m.manipulator] = &newManipulator
	}
	return newBrain
}

func randomChangeManipulatorValue(b *Brain) mChangeManipulatorValue {
	m := mChangeManipulatorValue{}
	size := len(b.ManipulationSystem)
	if size > 0 {
		m.manipulator = rand.Intn(size)
		m.direction = core.RandomDirection()
		m.delta = rand.Intn(manipulatorValueMaxDelta*2+1) - manipulatorValueMaxDelta
	}
	return m
}
