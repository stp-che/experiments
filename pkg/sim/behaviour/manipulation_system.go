package behaviour

import (
	"experiments/pkg/sim/core"
	"math/rand"
)

type Manipulator struct {
	ActionType ActionType
	DirValues  [8]int8
}

type ManipulationSystem []*Manipulator

func (s ManipulationSystem) ComputeIntention(activations ManipulationSystemActivation) *Intention {
	decisionTable := make(map[ActionType]*[8]int16)
	for direction, activation := range activations {
		dirCorretion := int(direction - core.Up)
		for i, power := range activation {
			mn := s[i]
			decisionRow, ok := decisionTable[mn.ActionType]
			if !ok {
				decisionRow = &[8]int16{}
				decisionTable[mn.ActionType] = decisionRow
			}
			for j, v := range mn.DirValues {
				decisionRow[(j+dirCorretion)%8] += int16(v) * power
			}
		}
	}
	var max int16 = 0
	var intention *Intention
	for aType, values := range decisionTable {
		for i, v := range values {
			if max < v {
				if intention == nil {
					intention = &Intention{}
				}
				intention.ActionType = aType
				intention.Direction = core.Direction(i + 1)
				max = v
			}
		}
	}
	return intention
}

func randomManipulationSystem() ManipulationSystem {
	res := make(ManipulationSystem, rand.Intn(5)+1)
	for i := 0; i < len(res); i++ {
		dirValues := [8]int8{}
		for i := 0; i < 8; i++ {
			dirValues[i] = int8(rand.Intn(255) - 128)
		}
		res[i] = &Manipulator{
			ActionType: randomActionType(),
			DirValues:  dirValues,
		}
	}
	return res
}
