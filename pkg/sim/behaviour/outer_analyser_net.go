package behaviour

import "math/rand"

type OuterAnalyzerLink struct {
	Analyzer    uint8
	Signal      uint8
	Manipulator uint8
	Power       int8
}

type OuterAnalyzerNet []*OuterAnalyzerLink

func (a OuterAnalyzerNet) Activation(signalTable CollectedOuterSignal) map[uint8]int16 {
	res := make(map[uint8]int16)
	for _, link := range a {
		if sig, present := signalTable[link.Analyzer]; present {
			if _, ok := res[link.Manipulator]; !ok {
				res[link.Manipulator] = 0
			}
			res[link.Manipulator] += int16(sig[link.Signal]) * int16(link.Power)
		}
	}
	return res
}

func randomOuterAnalyzerNet(analyzersCount, manipulatorsCount int) OuterAnalyzerNet {
	res := make(OuterAnalyzerNet, rand.Intn(analyzersCount)*rand.Intn(manipulatorsCount)+1)
	for i := 0; i < len(res); i++ {
		res[i] = randomOuterAnalyzerLink(analyzersCount, manipulatorsCount)
	}
	return res
}

func randomOuterAnalyzerLink(analyzersCount, manipulatorsCount int) *OuterAnalyzerLink {
	return &OuterAnalyzerLink{
		Analyzer:    uint8(rand.Intn(analyzersCount)),
		Signal:      uint8(rand.Intn(4)),
		Manipulator: uint8(rand.Intn(manipulatorsCount)),
		Power:       int8(rand.Intn(256) - 128),
	}
}
