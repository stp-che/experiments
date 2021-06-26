package behaviour

import "math/rand"

// {
//	 MinHealth, MaxHealth, OuterAnalyzerLink, Correction,
//   ...
// }
type HealthAnalyzerNet []uint8

const healthAnalyzerLinkSize = 4

func (n HealthAnalyzerNet) Correction(healthIndicator int) OuterAnalyzerNetCorrection {
	res := OuterAnalyzerNetCorrection{}
	for i := 0; i < len(n); i += healthAnalyzerLinkSize {
		minHealth, maxHealth, outerAnalyzerLink, correction := n[i], n[i+1], n[i+2], n[i+3]
		if minHealth > uint8(healthIndicator) || maxHealth < 255 && maxHealth < uint8(healthIndicator) {
			continue
		}

		if _, ok := res[outerAnalyzerLink]; !ok {
			res[outerAnalyzerLink] = 1.0
		}
		res[outerAnalyzerLink] *= float32(correction) / 16
	}
	return res
}

func (n HealthAnalyzerNet) normalize(bStruct BrainStructure) {
	for i := 0; i < int(bStruct.HealthAnalyzerNetSize); i++ {
		k := i * healthAnalyzerLinkSize
		if n[k+2] > bStruct.OuterAnalyzerNetSize {
			n[k+2] %= bStruct.OuterAnalyzerNetSize
		}
	}
}

func (n HealthAnalyzerNet) randomize(bStruct BrainStructure) {
	for i := 0; i < int(bStruct.HealthAnalyzerNetSize); i += healthAnalyzerLinkSize {
		minHealth := rand.Intn(200)
		n[i] = uint8(minHealth)
		n[i+1] = uint8(minHealth + rand.Intn(256-minHealth))
		n[i+2] = uint8(rand.Intn(int(bStruct.OuterAnalyzerNetSize)))
		n[i+3] = uint8(rand.Intn(256))
	}
}

// func randomHealthAnalyzerNet(outerAnalyzerNetSize int) HealthAnalyzerNet {
// 	count := rand.Intn((outerAnalyzerNetSize+1)/2) + 1
// 	res := make(HealthAnalyzerNet, count*healthAnalyzerLinkSize)
// 	for i := 0; i < count; i += healthAnalyzerLinkSize {
// 		minHealth := rand.Intn(200)
// 		res[i] = uint8(minHealth)
// 		res[i+1] = uint8(minHealth + rand.Intn(256-minHealth))
// 		res[i+2] = uint8(rand.Intn(outerAnalyzerNetSize))
// 		res[i+3] = uint8(rand.Intn(256))
// 	}
// 	return res
// }
