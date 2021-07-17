package behaviour

import (
	"math/rand"
)

const maxVisionRange = 15

type OuterReceptor struct {
	visionRange uint8
	cells       [][]uint8
}

func (r OuterReceptor) Size() int {
	return int(r.visionRange) * (int(r.visionRange) + 1)
}

func (r OuterReceptor) CollectSignal(signal []uint8) CollectedOuterSignal {
	res := make(CollectedOuterSignal)
	for i, sig := range signal {
		w := int(r.visionRange) + 1
		analyzer := r.cells[i/w][i%w]
		if _, ok := res[analyzer]; !ok {
			res[analyzer] = make(map[uint8]uint8)
		}
		res[analyzer][sig] += 1
	}
	return res
}

func (r OuterReceptor) copy() OuterReceptor {
	return r
}

func NewOuterReceptor(visionRange uint8) OuterReceptor {
	r := OuterReceptor{visionRange: visionRange}
	if visionRange == 0 {
		return r
	}

	r.cells = make([][]uint8, visionRange)
	for i := 0; i < int(visionRange); i++ {
		r.cells[i] = make([]uint8, visionRange+1)
	}
	return r
}

func randomOuterReceptor(outerAnalyzerSize int) OuterReceptor {
	n := rand.Intn(16) + 1
	res := NewOuterReceptor(uint8(n))
	for _, cellsRow := range res.cells {
		for i := 0; i < len(cellsRow); i++ {
			cellsRow[i] = uint8(rand.Intn(outerAnalyzerSize))
		}
	}
	return res
}
