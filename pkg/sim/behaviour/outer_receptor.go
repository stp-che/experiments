package behaviour

import (
	"fmt"
	"math/rand"
)

type OuterReceptor []uint8

func (r OuterReceptor) CollectSignal(signal []uint8) CollectedOuterSignal {
	res := make(CollectedOuterSignal)
	for i, sig := range signal {
		if i >= len(r) {
			fmt.Printf("%d\n%v\n%d\n%v\n", len(r), r, len(signal), signal)
		}
		analyzer := r[i]
		if _, ok := res[analyzer]; !ok {
			res[analyzer] = make([]uint8, 5)
		}
		res[analyzer][sig] += 1
	}
	return res
}

func randomOuterReceptor(outerAnalyzerSize int) OuterReceptor {
	n := rand.Intn(16) + 1
	res := make(OuterReceptor, n*(n+1))
	for i := 0; i < len(res); i++ {
		res[i] = uint8(rand.Intn(outerAnalyzerSize))
	}
	return res
}
