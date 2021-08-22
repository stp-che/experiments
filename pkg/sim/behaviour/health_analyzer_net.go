package behaviour

import (
	"fmt"
	"math/rand"
)

type HealthAnalyzerLink struct {
	MinHealth         uint8
	MaxHealth         uint8
	OuterAnalyzerLink uint8
	Correction        uint8
}

func (l HealthAnalyzerLink) copy() HealthAnalyzerLink {
	return l
}

const healthAnalyzerLinkCorrectionBase = 16

type HealthAnalyzerNet []*HealthAnalyzerLink

func (n HealthAnalyzerNet) String() string {
	s := "["
	for i, link := range n {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("&{%d, %d, %d, %d}", link.MinHealth, link.MaxHealth, link.OuterAnalyzerLink, link.Correction)
	}
	s += "]"

	return s
}

func (n HealthAnalyzerNet) Correction(healthIndicator int) OuterAnalyzerNetCorrection {
	res := OuterAnalyzerNetCorrection{}
	for _, link := range n {
		if link.MinHealth > uint8(healthIndicator) || link.MaxHealth < 255 && link.MaxHealth < uint8(healthIndicator) {
			continue
		}

		if _, ok := res[link.OuterAnalyzerLink]; !ok {
			res[link.OuterAnalyzerLink] = 1.0
		}
		res[link.OuterAnalyzerLink] *= float32(link.Correction) / healthAnalyzerLinkCorrectionBase
	}
	return res
}

func (n HealthAnalyzerNet) copy() HealthAnalyzerNet {
	nn := make(HealthAnalyzerNet, len(n))
	copy(nn, n)
	return nn
}

func (n HealthAnalyzerNet) appendSafely(link *HealthAnalyzerLink) HealthAnalyzerNet {
	size := len(n)
	nn := make(HealthAnalyzerNet, size+1)
	for i := 0; i < size; i++ {
		nn[i] = n[i]
	}
	nn[size] = link

	return nn
}

func randomHealthAnalyzerNet(outerAnalyzerNetSize int) HealthAnalyzerNet {
	count := rand.Intn((outerAnalyzerNetSize+1)/2) + 1
	res := make(HealthAnalyzerNet, count)
	for i := 0; i < count; i++ {
		res[i] = randomHealthAnalyzerLink(outerAnalyzerNetSize)
	}
	return res
}

func randomHealthAnalyzerLink(outerAnalyzerNetSize int) *HealthAnalyzerLink {
	minHealth := rand.Intn(200)
	return &HealthAnalyzerLink{
		MinHealth:         uint8(minHealth),
		MaxHealth:         uint8(minHealth + rand.Intn(256-minHealth)),
		OuterAnalyzerLink: uint8(rand.Intn(outerAnalyzerNetSize)),
		Correction:        uint8(rand.Intn(256)),
	}
}
