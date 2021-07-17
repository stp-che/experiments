package behaviour

import "math/rand"

type iMutation interface {
	Apply(*Brain) *Brain
}

func randomMutation(b *Brain) iMutation {
	switch rand.Intn(3) {
	case 0:
		return mIncreaseVisionRange{}
	case 1:
		return mDecreaseVisionRange{}
	case 2:
		return randomChangeOuterReceptor(b)
	default:
		return randomChangeOuterReceptor(b)
	}
}

type mIncreaseVisionRange struct{}

func (m mIncreaseVisionRange) Apply(b *Brain) *Brain {
	newBrain := b.copy()
	if b.OuterReceptor.visionRange == maxVisionRange {
		return newBrain
	}

	newRec := NewOuterReceptor(b.OuterReceptor.visionRange + 1)
	if b.OuterReceptor.visionRange > 0 {
		origCells := b.OuterReceptor.cells
		for i := 0; i < int(newRec.visionRange); i++ {
			for j := 0; j < int(newRec.visionRange)+1; j++ {
				ii, jj := i, j
				if ii > 0 {
					ii -= 1
				}
				if jj > 0 {
					jj -= 1
				}
				newRec.cells[i][j] = origCells[ii][jj]
			}
		}
	}

	newBrain.OuterReceptor = newRec

	return newBrain
}

type mDecreaseVisionRange struct{}

func (m mDecreaseVisionRange) Apply(b *Brain) *Brain {
	newBrain := b.copy()
	if b.OuterReceptor.visionRange == 0 {
		return newBrain
	}

	newRec := NewOuterReceptor(b.OuterReceptor.visionRange - 1)
	origCells := b.OuterReceptor.cells
	for i := 0; i < int(newRec.visionRange); i++ {
		for j := 0; j < int(newRec.visionRange)+1; j++ {
			newRec.cells[i][j] = origCells[i+1][j+1]
		}
	}

	newBrain.OuterReceptor = newRec

	return newBrain
}

type mChangeOuterReceptor struct {
	cell     int
	analyzer int
}

func (m mChangeOuterReceptor) Apply(b *Brain) *Brain {
	newBrain := b.copy()
	if m.cell >= b.OuterReceptor.Size() || m.analyzer >= b.OuterAnalyzersCount {
		return newBrain
	}

	w := int(b.OuterReceptor.visionRange) + 1
	mi := m.cell / w
	mj := m.cell % w

	newRec := NewOuterReceptor(b.OuterReceptor.visionRange)
	origCells := b.OuterReceptor.cells
	for i := 0; i < int(newRec.visionRange); i++ {
		for j := 0; j < int(newRec.visionRange)+1; j++ {
			if i == mi && j == mj {
				newRec.cells[i][j] = uint8(m.analyzer)
			} else {
				newRec.cells[i][j] = origCells[i][j]
			}
		}
	}

	newBrain.OuterReceptor = newRec

	return newBrain
}

func randomChangeOuterReceptor(b *Brain) mChangeOuterReceptor {
	m := mChangeOuterReceptor{}
	rSize := b.OuterReceptor.Size()
	if rSize > 0 && b.OuterAnalyzersCount > 0 {
		m.cell = rand.Intn(rSize)
		m.analyzer = rand.Intn(b.OuterAnalyzersCount)
	}

	return m
}
