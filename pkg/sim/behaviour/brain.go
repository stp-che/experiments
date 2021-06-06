package behaviour

import (
	"experiments/pkg/sim/core"
	"math/rand"
)

type OuterInput struct {
	Direction core.Direction
	Signal    []uint8
}

type InnerInput []int

type ProcessingResult struct {
	Decision   *Intention
	EnergyCost int
}

type IBrain interface {
	Process(OuterInput, InnerInput) *ProcessingResult
}

type manipulator struct {
	ActionType ActionType
	DirValues  [8]int8
}

type Brain struct {
	Manipulators []*manipulator
}

func (b *Brain) Process(o OuterInput, i InnerInput) *ProcessingResult {
	activity := make([]int16, len(b.Manipulators))
	for i := 0; i < len(activity); i++ {
		activity[i] = int16(rand.Intn(255) - 100)
	}
	return &ProcessingResult{
		Decision:   computeIntention(b.Manipulators, activity),
		EnergyCost: 10,
	}
}

func computeIntention(mnps []*manipulator, activity []int16) *Intention {
	decisionTable := make(map[ActionType]*[8]int16)
	// TODO: check genes and activity have the same len
	for i, mn := range mnps {
		decisionRow, ok := decisionTable[mn.ActionType]
		if !ok {
			decisionRow = &[8]int16{}
			decisionTable[mn.ActionType] = decisionRow
		}
		for j, v := range mn.DirValues {
			decisionRow[j] += int16(v) * activity[i]
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

func RandomBrain() *Brain {
	mnps := make([]*manipulator, rand.Intn(5)+1)
	for i := 0; i < len(mnps); i++ {
		dirValues := [8]int8{}
		for i := 0; i < 8; i++ {
			dirValues[i] = int8(rand.Intn(255) - 128)
		}
		mnps[i] = &manipulator{
			ActionType: randomActionType(),
			DirValues:  dirValues,
		}
	}
	return &Brain{
		Manipulators: mnps,
	}
}
