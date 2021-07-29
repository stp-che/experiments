package behaviour

import "experiments/pkg/sim/core"

const signalsCount = 4

type OuterInput struct {
	Direction core.Direction
	Signal    []uint8
}

type InnerInput []int

type ProcessingResult struct {
	Decision   *Intention
	EnergyCost int
}

// Signal being received by analyzer system in form {N: {i: Ai}}
// where N is an index of analyzer cell and Ai is a number of signal i received by the analyzer cell
// TODO: think about converting to map[uint8]map[uint8]uint8
type CollectedOuterSignal map[uint8]map[uint8]uint8

type OuterAnalyzerNetCorrection map[uint8]float32

// Activation values for manipulation system in form {direction: {i: pow}}
// where direction - from where outer signal was collected
// i - number of manipulator,
// pow - power in which manipulator i is activated
type ManipulationSystemActivation map[core.Direction]map[uint8]int16

type iMutation interface {
	apply(*Brain) *Brain
}

func cutToByte(i int) uint8 {
	if i < 0 {
		return 0
	}
	if i > 255 {
		return 255
	}
	return uint8(i)
}

func cutToInt8(i int) int8 {
	if i < -128 {
		return -128
	}
	if i > 127 {
		return 127
	}
	return int8(i)
}
