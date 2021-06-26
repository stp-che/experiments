package behaviour

import (
	"experiments/pkg/sim/core"
	"math/rand"
)

type ManipulationSystem []uint8

const manipulatorSize = 9

func (s ManipulationSystem) ComputeIntention(activations ManipulationSystemActivation) *Intention {
	decisionTable := make(map[ActionType]*[8]int16)
	for direction, activation := range activations {
		dirCorretion := int(direction - core.Up)

		for manipulatorIndex, power := range activation {
			i := int(manipulatorIndex) * manipulatorSize
			aType := actionType(s[i])
			decisionRow, ok := decisionTable[aType]
			if !ok {
				decisionRow = &[8]int16{}
				decisionTable[aType] = decisionRow
			}
			for j := 0; j < 8; j++ {
				decisionRow[(j+dirCorretion)%8] += (int16(s[i+j+1]) - 128) * power
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

func (s ManipulationSystem) normalize(bStruct BrainStructure) {
	for mi := 0; mi < int(bStruct.ManipulationSystemSize); mi++ {
		i := mi * manipulatorSize
		s[i] = uint8(actionType(s[i]))
	}
}

func (s ManipulationSystem) randomize(bStruct BrainStructure) {
	for i := 0; i < int(bStruct.ManipulationSystemSize); i++ {
		j := i * manipulatorSize
		s[j] = uint8(randomActionType())
		for k := 0; k < 8; k++ {
			s[j+k] = uint8(rand.Intn(256))
		}
	}
}
