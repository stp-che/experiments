package behaviour

import "math/rand"

const (
	healthAnalyzerCorrectionMaxDelta = 10
	healthAnalyzerRangeMaxDelta      = 10
)

type mAddHealthAnalyzerLink struct {
	link HealthAnalyzerLink
}

func (m mAddHealthAnalyzerLink) apply(brain *Brain) *Brain {
	newBrain := brain.copy()
	if int(m.link.OuterAnalyzerLink) < len(brain.OuterAnalyzerNet) {
		newBrain.HealthAnalyzerNet = append(brain.HealthAnalyzerNet, &m.link)
	}
	return newBrain
}

func randomAddHealthAnalyzerLink(b *Brain) mAddHealthAnalyzerLink {
	m := mAddHealthAnalyzerLink{}
	size := len(b.OuterAnalyzerNet)
	if size > 0 {
		m.link = *randomHealthAnalyzerLink(size)
	}
	return m
}

type mChangeHealthAnalyzerCorrection struct {
	healthAnalyzerLink uint8
	delta              int
}

func (m mChangeHealthAnalyzerCorrection) apply(brain *Brain) *Brain {
	newBrain := brain.copy()
	if int(m.healthAnalyzerLink) < len(brain.HealthAnalyzerNet) && m.delta != 0 {
		c := int(newBrain.HealthAnalyzerNet[m.healthAnalyzerLink].Correction) + m.delta
		if c == healthAnalyzerLinkCorrectionBase {
			newBrain.HealthAnalyzerNet = removeHealthAnalyserLink(newBrain.HealthAnalyzerNet, int(m.healthAnalyzerLink))
		} else {
			newBrain.HealthAnalyzerNet = changeHealthAnalyserLink(
				newBrain.HealthAnalyzerNet,
				int(m.healthAnalyzerLink),
				func(link *HealthAnalyzerLink) { link.Correction = cutToByte(c) },
			)
		}
	}
	return newBrain
}

func randomChangeHealthAnalyzerCorrection(b *Brain) mChangeHealthAnalyzerCorrection {
	m := mChangeHealthAnalyzerCorrection{}
	if size := len(b.HealthAnalyzerNet); size > 0 {
		m.healthAnalyzerLink = uint8(rand.Intn(size))
		m.delta = rand.Intn(healthAnalyzerCorrectionMaxDelta*2+1) - healthAnalyzerCorrectionMaxDelta
	}
	return m
}

type mChangeHealthAnalyzerMinMax struct {
	healthAnalyzerLink uint8
	deltaMin           int
	deltaMax           int
}

func (m mChangeHealthAnalyzerMinMax) apply(brain *Brain) *Brain {
	newBrain := brain.copy()
	if int(m.healthAnalyzerLink) < len(brain.HealthAnalyzerNet) && m.deltaMin != 0 && m.deltaMax != 0 {
		targetLink := newBrain.HealthAnalyzerNet[m.healthAnalyzerLink]
		newMin := cutToByte(int(targetLink.MinHealth) + m.deltaMin)
		newMax := cutToByte(int(targetLink.MaxHealth) + m.deltaMax)
		if newMin >= newMax {
			newBrain.HealthAnalyzerNet = removeHealthAnalyserLink(newBrain.HealthAnalyzerNet, int(m.healthAnalyzerLink))
		} else {
			newBrain.HealthAnalyzerNet = changeHealthAnalyserLink(
				newBrain.HealthAnalyzerNet,
				int(m.healthAnalyzerLink),
				func(link *HealthAnalyzerLink) {
					link.MinHealth = newMin
					link.MaxHealth = newMax
				},
			)
		}
	}

	return newBrain
}

func randomChangeHealthAnalyzerMinMax(b *Brain) mChangeHealthAnalyzerMinMax {
	m := mChangeHealthAnalyzerMinMax{}
	if size := len(b.HealthAnalyzerNet); size > 0 {
		m.healthAnalyzerLink = uint8(rand.Intn(size))
		delta := rand.Intn(healthAnalyzerRangeMaxDelta*2+1) - healthAnalyzerRangeMaxDelta
		if rand.Intn(2) == 0 {
			m.deltaMin = delta
		} else {
			m.deltaMax = delta
		}
	}
	return m
}

func removeHealthAnalyserLink(net HealthAnalyzerNet, idx int) HealthAnalyzerNet {
	newNet := make(HealthAnalyzerNet, len(net)-1)
	for i, link := range net {
		if i < idx {
			newNet[i] = link
		} else if i > idx {
			newNet[i-1] = link
		}
	}
	return newNet
}

func changeHealthAnalyserLink(net HealthAnalyzerNet, idx int, change func(*HealthAnalyzerLink)) HealthAnalyzerNet {
	newNet := net.copy()
	newLink := newNet[idx].copy()
	change(&newLink)
	newNet[idx] = &newLink
	return newNet
}
