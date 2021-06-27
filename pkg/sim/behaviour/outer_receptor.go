package behaviour

import (
	"math/rand"
)

type OuterReceptor []uint8

func (r OuterReceptor) CollectSignal(signal []uint8) CollectedOuterSignal {
	res := make(CollectedOuterSignal)
	for i, sig := range signal {
		analyzer := r[i]
		if _, ok := res[analyzer]; !ok {
			res[analyzer] = make(map[uint8]uint8)
		}
		res[analyzer][sig] += 1
	}
	return res
}

func (r OuterReceptor) normalize(s BrainStructure) {
	for i := 0; i < len(r); i++ {
		if r[i] >= s.OuterAnalyzersCount {
			r[i] %= s.OuterAnalyzersCount
		}
	}
}

func (r OuterReceptor) randomize(s BrainStructure) {
	for i := 0; i < len(r); i++ {
		r[i] = uint8(rand.Intn(int(s.OuterAnalyzersCount)))
	}
}
