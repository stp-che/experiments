package behaviour

import "math/rand"

const outerAnalyzerPowerMaxDelta = 10

type mChangeOuterAnalyzerLinkPower struct {
	analyzer    uint8
	signal      uint8
	manipulator uint8
	delta       int
}

func (m mChangeOuterAnalyzerLinkPower) apply(brain *Brain) *Brain {
	newBrain := brain.copy()
	i := m.linkIndex(brain)
	if i == -1 {
		if m.analyzer < uint8(brain.OuterAnalyzersCount) && m.manipulator < uint8(len(brain.ManipulationSystem)) {
			newBrain.OuterAnalyzerNet = append(newBrain.OuterAnalyzerNet, &OuterAnalyzerLink{
				Analyzer:    m.analyzer,
				Signal:      m.signal,
				Manipulator: m.manipulator,
				Power:       cutToInt8(m.delta),
			})
		}
	} else {
		link := newBrain.OuterAnalyzerNet[i]
		newPower := cutToInt8(int(link.Power) + m.delta)
		if newPower == 0 {
			if !m.linkHasIncomingRefs(brain, uint8(i)) {
				newBrain.OuterAnalyzerNet = removeOuterAnalyzerLink(newBrain.OuterAnalyzerNet, i)
			}
		} else {
			newBrain.OuterAnalyzerNet = newBrain.OuterAnalyzerNet.copy()
			newLink := newBrain.OuterAnalyzerNet[i].copy()
			newLink.Power = newPower
			newBrain.OuterAnalyzerNet[i] = &newLink
		}
	}
	return newBrain
}

func (m mChangeOuterAnalyzerLinkPower) linkIndex(brain *Brain) int {
	for i, link := range brain.OuterAnalyzerNet {
		if link.Analyzer == m.analyzer && link.Signal == m.signal && link.Manipulator == m.manipulator {
			return i
		}
	}
	return -1
}

func (m mChangeOuterAnalyzerLinkPower) linkHasIncomingRefs(brain *Brain, i uint8) bool {
	for _, link := range brain.HealthAnalyzerNet {
		if link.OuterAnalyzerLink == i {
			return true
		}
	}
	return false

}

func randomChangeOuterAnalyzerLinkPower(b *Brain) mChangeOuterAnalyzerLinkPower {
	m := mChangeOuterAnalyzerLinkPower{}
	if b.OuterAnalyzersCount > 0 && len(b.ManipulationSystem) > 0 {
		m.analyzer = uint8(rand.Intn(b.OuterAnalyzersCount))
		m.signal = uint8(rand.Intn(signalsCount))
		m.manipulator = uint8(rand.Intn(len(b.ManipulationSystem)))
		m.delta = rand.Intn(outerAnalyzerPowerMaxDelta*2+1) - outerAnalyzerPowerMaxDelta
	}
	return m
}

func removeOuterAnalyzerLink(net OuterAnalyzerNet, idx int) OuterAnalyzerNet {
	newNet := make(OuterAnalyzerNet, len(net)-1)
	for i, link := range net {
		if i < idx {
			newNet[i] = link
		} else if i > idx {
			newNet[i-1] = link
		}
	}
	return newNet
}
