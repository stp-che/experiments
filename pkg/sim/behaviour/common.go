package behaviour

import (
	"experiments/pkg/sim/core"
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

// Signal being received by analyzer system in form {N: {A0, A1, ..., An}}
// where N is an index of analyzer cell and Ai is a number of signal i received by the analyzer cell
// TODO: think about converting to map[uint8]map[uint8]uint8
type CollectedOuterSignal map[uint8][]uint8

type OuterAnalyzerNetCorrection map[uint8]float32

// Activation values for manipulation system in form {direction: {i: pow}}
// where direction - from where outer signal was collected
// i - number of manipulator,
// pow - power in which manipulator i is activated
type ManipulationSystemActivation map[core.Direction]map[uint8]int16
