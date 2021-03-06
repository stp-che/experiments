package behaviour

import (
	"reflect"
	"testing"
)

func TestActivation(t *testing.T) {
	testNet := OuterAnalyzerNet{
		// Analyzer, Signal, Manipulator, Power
		0, 0, 5, 138,
		0, 0, 8, 228,
		0, 1, 5, 123,
		1, 0, 8, 178,
	}

	cases := []struct {
		Net        OuterAnalyzerNet
		Signal     CollectedOuterSignal
		Correction OuterAnalyzerNetCorrection
		Result     map[uint8]int16
	}{
		{
			Net: testNet,
			Signal: CollectedOuterSignal{
				0: {0: 2, 1: 1},
				1: {0: 1, 1: 1},
			},
			Correction: OuterAnalyzerNetCorrection{},
			Result:     map[uint8]int16{5: 15, 8: 250},
		},
		{
			Net: testNet,
			Signal: CollectedOuterSignal{
				0: {0: 2, 1: 1},
				1: {0: 1, 1: 1},
			},
			Correction: OuterAnalyzerNetCorrection{1: 0.11},
			Result:     map[uint8]int16{5: 15, 8: 72},
		},
		{
			Net:        testNet,
			Signal:     CollectedOuterSignal{2: {0: 2, 1: 1}},
			Correction: OuterAnalyzerNetCorrection{},
			Result:     map[uint8]int16{},
		},
		{
			Net:        testNet,
			Signal:     CollectedOuterSignal{0: {99: 2}},
			Correction: OuterAnalyzerNetCorrection{},
			Result:     map[uint8]int16{},
		},
	}

	for i, c := range cases {
		res := c.Net.Activation(c.Signal, c.Correction)
		if !reflect.DeepEqual(res, c.Result) {
			t.Errorf("[%d] expect result to eq %v, got %v", i, c.Result, res)
		}
	}
}
