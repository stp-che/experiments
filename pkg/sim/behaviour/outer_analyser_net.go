package behaviour

import (
	"math/rand"
)

type OuterAnalyzerLink struct {
	Analyzer    uint8
	Signal      uint8
	Manipulator uint8
	Power       int8
}

// {
//	 Analyzer, Signal, Manipulator, Power,
//   ...
// }
type OuterAnalyzerNet []uint8

const outerAnalyzerLinkSize = 4

func (a OuterAnalyzerNet) Activation(signalTable CollectedOuterSignal, correction OuterAnalyzerNetCorrection) map[uint8]int16 {
	res := make(map[uint8]int16)
	count := len(a) / outerAnalyzerLinkSize
	for link := 0; link < count; link++ {
		i := link * outerAnalyzerLinkSize
		analyzer, signal, manipulator, power := a[i], a[i+1], a[i+2], int16(a[i+3])-128
		if sig, present := signalTable[analyzer]; present {
			if value, signalPresent := sig[signal]; signalPresent {
				influence := int16(value) * int16(power)
				if corr, ok := correction[uint8(link)]; ok {
					influence = int16(float32(influence) * corr)
				}
				res[manipulator] += influence
			}
		}
	}
	return res
}

func (a OuterAnalyzerNet) normalize(s BrainStructure) {
	count := len(a) / outerAnalyzerLinkSize
	for link := 0; link < count; link++ {
		i := link * outerAnalyzerLinkSize
		if a[i] >= s.OuterAnalyzersCount {
			a[i] %= s.OuterAnalyzersCount
		}
		if a[i+2] >= s.ManipulationSystemSize {
			a[i+2] %= s.ManipulationSystemSize
		}
	}
}

func (a OuterAnalyzerNet) randomize(bStruct BrainStructure) {
	for i := 0; i < int(bStruct.OuterAnalyzerNetSize); i++ {
		j := i * outerAnalyzerLinkSize
		a[j] = uint8(rand.Intn(int(bStruct.OuterAnalyzersCount)))
		a[j+1] = uint8(rand.Intn(4))
		a[j+2] = uint8(rand.Intn(int(bStruct.ManipulationSystemSize)))
		a[j+3] = uint8(rand.Intn(256))
	}
}

// func randomOuterAnalyzerNet(analyzersCount, manipulatorsCount int) OuterAnalyzerNet {
// 	count := rand.Intn(analyzersCount)*rand.Intn(manipulatorsCount) + 1
// 	res := make(OuterAnalyzerNet, count*outerAnalyzerLinkSize)
// 	for i := 0; i < len(res); i += outerAnalyzerLinkSize {
// 		res[i] = uint8(rand.Intn(analyzersCount))
// 		res[i+1] = uint8(rand.Intn(4))
// 		res[i+2] = uint8(rand.Intn(manipulatorsCount))
// 		res[i+3] = uint8(rand.Intn(256))
// 	}
// 	return res
// }
